package main

import (
	"fmt"
	"runtime/debug"

	"github.com/freeconf/yang/fc"
	"github.com/manojankitha/redis-proxy/pkg/config"
	"github.com/manojankitha/redis-proxy/pkg/http/proxy"
)

// System Description: Redis Proxy
// Author: Ankitha Manoj

func main() {
	defer func() {
		if r := recover(); r != nil {
			fc.Err.Printf("Main error: %s\nFailed to create Redis Proxy Service %s", r, debug.Stack())
		}
	}()

	fmt.Println("Starting Redis proxy server on port 9000...")
	config := config.LoadConfig()

	proxyService, err := proxy.NewRedisProxy(config)
	if err != nil {
		panic("Failed to create Redis Proxy Service")
	}
	fc.Info.Printf("Starting Redis Proxy Server on port %d ...", config.ProxyPort)
	server := proxy.NewRedisProxyService(proxyService, config.ProxyPort)
	server.ListenAndServe()

}
