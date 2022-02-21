# redis-proxy

 A HTTP based Redis Proxy Service that caches for GET requests

## High-level architecture overview

In a standalone redis caching system, each instance of an application connects to the redis instance/cluster. As the application scales (#instances increase or rate of requests increases) the overhead on the redis server increases reducing the performance of our application. This can be resolved by having a proxy service that handles the requests before it hits the redis layer.

<img src=images/systemOverview.jpg Level width="70%" title="System Overview">

Each instance of the proxy service is associated with a single Redis service instance called the “backing Redis”. Clients interface to the Redis proxy through HTTP, with the Redis “GET” command mapped to the HTTP “GET” method. A GET request, directed at the proxy, returns the value of the specified key from the proxy’s local cache if the local cache contains a value for that key. If the local cache does not contain a value for the specified key, it fetches the value from the backing Redis instance, using the Redis GET command, and stores it in the local cache, associated with the specified key.

Assumptions:

1. The proxy is a simple read through cache that only accepts GET requests.
2. All writes directed directly to the backing Redis instance, treating the proxy out-of-band and ignoring it.
3. The cache stores only string data types.

Once the cache fills to capacity, the least recently used (i.e.read) key is evicted each time a new key needs to be added to the cache.
To implement Least Recently Used (LRU) eviction policy for the cache, I have imported an external library `https://github.com/hashicorp/golang-lru/blob/master/simplelru/lru.go` that uses a Map and a Doubly-linked List for storing everything. Both map and doubly-linked list are used so that get and set operations are of O(1) even with evictions.

## Code Brief

### Folder structure

```
redis-proxy
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── redis-proxy
│       ├── coverage.out
│       └── main.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── images
│   └── systemOverview.jpg
└── pkg
    ├── cache
    │   ├── cache.go
    │   ├── cache_test.go
    │   └── coverage.out
    ├── config
    │   └── config.go
    ├── http
    │   └── proxy
    │       ├── coverage.out
    │       ├── handler.go
    │       ├── proxy.go
    │       └── proxy_test.go
    └── redis
        ├── coverage.out
        ├── redis.go
        └── redis_test.go
```

The project contains of two main folders `cmd` and `pkg`. Overview of the packages:

### `config/config.go`

The `Config` struct stores the configurable parameters of a redis proxy service. These are the configurable parameters that can be set in an environment file.

```
REDIS_ADDRESS=#redis address host:port ex: "redis:6379"
GLOBAL_CACHE_EXPIRY_TIME=#local cache expiry time in milliseconds
CACHE_CAPACITY=#cache capacity in int
PROXY_PORT=#port proxy is running on
MAX_CLIENTS=#concurrent client limit
```

### `redis/redis.go`

This is a wrapper code around `go-redis` redis client. It handles initialization of redis client, `SetVal`, `GetVal` and ensures redis storage DB is flushed afterwards. This package helps create the backing redis server instance that all proxy instances connect to.

### `cache/cache.go`

The cache package initializes the cache, handles `Setval`, `GetVal` and ensures cache is purged after use.

Values added to the local cache are expired after being in the cache for a time duration that is globally configured per instance using the environment variable `GLOBAL_CACHE_EXPIRY_TIME` (in ms). Every Get request to the cache first checks if the value associated with the key has expired. If a value expires, a GET request to it will act as if the value associated with the key has never been stored in the cache. If the value has not expired, it returns the value.

The cache capacity is configured using the environment variable `CACHE_CAPACITY`.

The cache supports sequential concurrent processing where in multiple clients are able to concurrently connect to the proxy without adversely impacting the functional behavior of the proxy. To set or get keys, any access to the cache must acquire a mutex. the `LocalCacheClient` struct stores the pointer to the LRU cache and the mutex

### `http/proxy/proxy.go`

The `RedisProxy` struct stores `RedisStorage` -  Redis client, `proxyStorage` - a local cache client, maxClients - to specify max limit on concurrent clients that the proxy service can support ensuring consistency, `sema` - a semaphore to support parallel concurrent processing.

`NewRedisProxy(config *Config)` returns a new redis proxy for the given configuration.
`NewRedisProxyService(proxy *RedisProxy, port int)` returns a HTTP server running on port specified through which Clients interface to the Redis proxy server.

### `http/proxy/handler.go`

This redirects HTTP requests to the appropriate handlers functions.
`GET /` - displays help menu for redis proxy service
`GET /GET/{keyZ}` - returns the value associated with `keyZ`.

A GET request, directed at the proxy, returns the value of the
specified key from the proxy’s local cache if the local cache contains a value for that key. If the local cache does not contain
a value for the specified key, it fetches the value from the backing Redis instance, using the Redis GET command, and
stores it in the local cache, associated with the specified key.

### `cmd/main.go`

This is the entrypoint for the project. Based on the configuration passed by the user, it initializes a new redis proxy server that clients can interface with using a HTTP web service.

## Instructions for how to run the proxy and tests

 Prerequisites:
 `make`, `docker`, `docker-compose`

 Run from Source:

 ```
 # clone repo
 git clone git@github.com:manojankitha/redis-proxy.git

 cd redis-proxy

 # Run the proxy:
 ./redis-proxy

 # Unit test:
 make test

 # Single Click build and test:
 make build test
```

 Run from container:

 ```
 docker-compose up -d --build
 ```

## Test coverage

Run the following command to calculate the coverage for your current unit test:

`go test ./... -coverprofile=coverage.out`

Go saves this coverage data in the file coverage.out. Now you can present the results in a web browser.
Run the following command:

`go tool cover -html=coverage.out`

A web browser will open, and your results will appear in the browser.

## Algorithmic Complexity

Redis: `GetVal` and `SetVal` operations would take `O(1)` time considering our application only considers storing string datatype.
Note: For other datatypes depending on what is stored, it may vary (Ref: `https://github.com/ZhenningLang/redis-command-complexity-cheatsheet`)

Local cache: The cache data structure is a simple map. Every `GetVal` and `SetVal` takes `O(1)` time.
Deletion is based on eviction algorithm - LRU which also takes O(1) time (Ref: `https://github.com/hashicorp/golang-lru/blob/master/simplelru/lru.go`)

Overall time complexity for API `GET/{key}`: `O(1)`

## Time Spent - 16 hours

1. Initial set up of HTTP server - 1 hour
2. Configuration parsing and testing -  1 hour
3. Redis service and testing - 2 hours
4. Local Cache Client and testing - 3 hours
5. Redis Proxy service and testing - 3 hours
6. Parallel Concurrent processing - 1 hours
7. Dockerize - 3 hours (had to learn and do)
8. Documentation - 2 hours

## Future Improvements

1. Add unit and regression testing for parallel concurrent processing across multiple clients
