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
