package proxy

import (
	"net"
	"net/http"
	"strconv"
	"testing"

	"github.com/manojankitha/redis-proxy/pkg/config"
	"github.com/stretchr/testify/assert"
)

// define test config
const TEST_CACHE_CAPACITY = 5
const TEST_GLOBAL_CACHE_EXPIRY_TIME = 120000 // in ms
const TEST_REDIS_ADDRESS = "localhost:6379"
const TEST_PROXY_PORT = 8080

// Test to initialize new redis proxy
func TestRedisProxyInitialization(t *testing.T) {
	testConfig := &config.Config{
		RedisAddr:             TEST_REDIS_ADDRESS,
		GlobalCacheExpiryTime: TEST_GLOBAL_CACHE_EXPIRY_TIME,
		CacheCapacity:         TEST_CACHE_CAPACITY,
		ProxyPort:             TEST_PROXY_PORT,
	}

	_, err := NewRedisProxy(testConfig)
	assert.NoError(t, err)
}

// Test to confirm redis proxy service running on a specific port
func TestRedisProxyService(t *testing.T) {
	testConfig := &config.Config{
		RedisAddr:             TEST_REDIS_ADDRESS,
		GlobalCacheExpiryTime: TEST_GLOBAL_CACHE_EXPIRY_TIME,
		CacheCapacity:         TEST_CACHE_CAPACITY,
		ProxyPort:             TEST_PROXY_PORT,
	}
	proxyService, err := NewRedisProxy(testConfig)
	if err != nil {
		panic("Failed to create Redis Proxy Service")
	}
	server := NewRedisProxyService(proxyService, testConfig.ProxyPort)
	l, err := net.Listen("tcp", ":"+strconv.Itoa(testConfig.ProxyPort))
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	go server.Serve(l)
	host := l.Addr().String()

	res, err := http.Get("http://" + host + "/")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode, "Redis Proxy Server Failed to Initialize")

	res, err = http.Get("http://" + host + "/GET")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusNotFound, res.StatusCode)

	res, err = http.Get("http://" + host + "/GET/{key}")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, res.StatusCode, "Redis Proxy Server Failed to Initialize")

}
