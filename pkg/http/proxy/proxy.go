package proxy

import (
	"net/http"
	"strconv"

	"github.com/freeconf/yang/fc"
	"github.com/manojankitha/redis-proxy/pkg/cache"
	"github.com/manojankitha/redis-proxy/pkg/config"
	"github.com/manojankitha/redis-proxy/pkg/redis"
)

type RedisProxy struct {
	redisStorage *redis.RedisClient
	proxyStorage *cache.LocalCacheClient
	sema         chan struct{} // reference: https://www.pauladamsmith.com/blog/2016/04/max-clients-go-net-http.html
	maxClients   int
}

func NewRedisProxy(pConfig *config.Config) (*RedisProxy, error) {
	redisServ, err := redis.RedisInit(pConfig.RedisAddr)
	if err != nil {
		return nil, err
	}
	proxyServ, err := cache.CacheInit(pConfig.CacheCapacity, pConfig.GlobalCacheExpiryTime)
	if err != nil {
		return nil, err
	}
	fc.Info.Println("Initializing redis proxy service... ")
	service := &RedisProxy{
		redisStorage: redisServ,
		proxyStorage: proxyServ,
		sema:         make(chan struct{}, pConfig.MaxClients),
		maxClients:   pConfig.MaxClients,
	}
	return service, nil
}

func NewRedisProxyService(proxy *RedisProxy, port int) *http.Server {

	router := InitHandlers(proxy)
	server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: router,
	}
	return server

}
