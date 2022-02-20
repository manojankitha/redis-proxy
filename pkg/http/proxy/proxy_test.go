package proxy

import (
	"io/ioutil"
	"net"
	"net/http"
	"runtime/debug"
	"strconv"
	"testing"

	"github.com/freeconf/yang/fc"
	"github.com/manojankitha/redis-proxy/pkg/config"
	"github.com/stretchr/testify/assert"
)

// define test config
const TEST_CACHE_CAPACITY = 5
const TEST_GLOBAL_CACHE_EXPIRY_TIME = 120000 // in ms
const TEST_REDIS_ADDRESS = "localhost:6379"
const TEST_PROXY_PORT = 9000
const MAX_CLIENTS = 2

// Test to initialize new redis proxy
func TestRedisProxyInitialization(t *testing.T) {
	testConfig := &config.Config{
		RedisAddr:             TEST_REDIS_ADDRESS,
		GlobalCacheExpiryTime: TEST_GLOBAL_CACHE_EXPIRY_TIME,
		CacheCapacity:         TEST_CACHE_CAPACITY,
		ProxyPort:             TEST_PROXY_PORT,
		MaxClients:            MAX_CLIENTS,
	}

	_, err := NewRedisProxy(testConfig)
	assert.NoError(t, err)
}

// Make request to specified path and return
// Copied from: https://github.com/segmentio/testdemo/blob/master/web/web_test.go#L13
func getRequestBody(path string, s *http.Server, t *testing.T) string {
	res := getResponse(path, s, t)

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	return string(body)
}

func getResponse(path string, s *http.Server, t *testing.T) *http.Response {
	// Pick port automatically for parallel tests and to avoid conflicts
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	go s.Serve(l)

	res, err := http.Get("http://" + l.Addr().String() + path)
	if err != nil {
		t.Fatal(err)
	}

	return res
}

// Test to confirm redis proxy service running on a specific port and is working as per requirements
func TestRedisProxyService(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			fc.Err.Printf("TestRedisProxyService error: %s\nFailed to create Redis Proxy Service %s", r, debug.Stack())
		}
	}()
	testConfig := &config.Config{
		RedisAddr:             TEST_REDIS_ADDRESS,
		GlobalCacheExpiryTime: TEST_GLOBAL_CACHE_EXPIRY_TIME,
		CacheCapacity:         TEST_CACHE_CAPACITY,
		ProxyPort:             TEST_PROXY_PORT,
		MaxClients:            MAX_CLIENTS,
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

	proxyService.redisStorage.SetKey("key1", "valueSetInRedis", 0)

	response := getResponse("/", server, t)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	response = getResponse("/GET", server, t)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
	response = getResponse("/GET/key", server, t)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Request value from proxy
	req := getRequestBody("/GET/key1", server, t)

	assert.Equal(t, "valueSetInRedis", string(req), "proxy incorrectly retrieved value from redis")

	// Update key in redis
	proxyService.proxyStorage.SetKey("key1", "valueModifiedByProxy")

	req = getRequestBody("/GET/key1", server, t)
	assert.Equal(t, "valueModifiedByProxy", string(req), "proxy incorrectly retrieved value from cache")

	assert.Equal(t, http.StatusOK, res.StatusCode, "Redis Proxy Server Failed to Initialize")
	proxyService.proxyStorage.PurgeCache()
	proxyService.redisStorage.FlushRedisDB()

}
