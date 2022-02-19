package proxy

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/freeconf/yang/fc"
	"github.com/gorilla/mux"
)

func InitHandlers(ps *RedisProxy) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/", ps.helpMenuHandler()).Methods("GET")
	router.HandleFunc("/GET/{key}", ps.getRequestHandler()).Methods("GET")
	return router
}

func (ps *RedisProxy) helpMenuHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "**************************************************************************** ")
		fmt.Fprintln(w, "*                                                                          * ")
		fmt.Fprintln(w, "*                                                                          * ")
		fmt.Fprintln(w, "*                 WELCOME TO HELP MENU FOR REDIS PROXY                     * ")
		fmt.Fprintln(w, "*                                                                          * ")
		fmt.Fprintln(w, "*                                                                          * ")
		fmt.Fprintln(w, "**************************************************************************** ")
		fmt.Fprintln(w, "API Instructions:")
		fmt.Fprintln(w, "GET /GET/{key} -  returns value of specified key from proxy's local cache if the local cache contains a value for that key. If the local cache does not contain a value for the specified key, it fetches the value from the backing Redis instance.")
	}
}

func (ps *RedisProxy) getRequestHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				fc.Err.Printf("GetRequestHandler error: %s\n Error details: %s", r, debug.Stack())
			}
		}()

		key := mux.Vars(r)["key"]

		// check if key present in local cache
		if proxyCacheValue, exists := ps.proxyStorage.GetKey(key); exists {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(proxyCacheValue))
			return
		}

		// get value from redis
		redisValue, err := ps.redisStorage.GetKey(key)

		if err == nil {
			ps.redisStorage.SetKey(key, redisValue, time.Minute*100) // making an assumption here
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(redisValue))
			return

		} else {
			log.Printf("Error retrieving value of %s from Redis. Error details: %s", key, err)
			// if redis service is closed or connection is lost
			if err := ps.redisStorage.C.Ping().Err(); err != nil {
				panic("Failure connecting to redis service. Suspecting redis service is stopped. " + err.Error())
			}

		}

	}
}
