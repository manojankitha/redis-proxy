package config

import (
	"flag"
)

type Config struct {
	RedisAddr             string `arg:"env:REDIS_ADDRESS"`
	GlobalCacheExpiryTime int    `arg:"env:CACHE_EXPIRY_TIME"`
	CacheCapacity         int    `arg:"env:CACHE_CAPACITY"`
	ProxyPort             int    `arg:"env:PROXY_PORT"`
}

func LoadConfig() *Config {
	config := &Config{}
	flag.StringVar(&config.RedisAddr, "redis-addr", "localhost:6379", "Address of backing Redis server")                          // using default redis address port
	flag.IntVar(&config.GlobalCacheExpiryTime, "cache-expiry-time", 60*1000, "Cache expiry time(Please mention in milliseconds)") // todo for future: Handle time input as duration because user should not concern themselves with conversion
	flag.IntVar(&config.CacheCapacity, "cache-capacity", 500, "Cache capacity")
	flag.IntVar(&config.ProxyPort, "proxy-port", 8080, "Port the proxy server listens on")
	flag.Parse()
	return config

}
