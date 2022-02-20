package cache

import (
	"log"
	"sync"
	"time"

	"github.com/freeconf/yang/fc"
	lru "github.com/hashicorp/golang-lru"
)

type cacheEntry struct {
	value  string // assuming all data can be stored as string type
	expiry int64
}

type LocalCacheClient struct {
	Cache           *lru.Cache
	cacheExpiryTime int64
	mutex           *sync.Mutex
}

func CacheInit(cacheCapacity int, cacheExpiry int) (*LocalCacheClient, error) {
	c, err := lru.New(cacheCapacity)
	if err != nil {
		log.Printf("Failed to create proxy cache: %v.", err)
		return nil, err
		//panic("Unable to create proxy cache: " + err.Error()) // since the assigned task is to build proxy client, it is crucial to create proxy cache
	}
	fc.Info.Println("Initializing proxy cache...")
	cacheClient := &LocalCacheClient{
		Cache:           c,
		cacheExpiryTime: int64(cacheExpiry),
		mutex:           &sync.Mutex{},
	}

	return cacheClient, nil
}

//GetKey get key
func (c *LocalCacheClient) GetKey(key string) (string, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if val, ok := c.Cache.Get(key); ok {
		currValue := val.(*cacheEntry)

		if !c.hasExpired(currValue) {
			return currValue.value, true
		}
		c.Cache.Remove(key)
		return "", false

	}
	return "", false

}

// has key expired fn
func (c *LocalCacheClient) hasExpired(ce *cacheEntry) bool {
	now := time.Now().UnixNano() / 1e6 // now time in milliseconds
	// if cache persistance time exceeded expiry time(inclusive limits)
	return now >= ce.expiry
}

//SetKey set key
func (c *LocalCacheClient) SetKey(key string, value string) error {
	now := time.Now().UnixNano() / 1e6 // now time in milliseconds
	expiry := now + c.cacheExpiryTime

	cacheEntry := &cacheEntry{value: value, expiry: expiry}
	c.mutex.Lock()
	c.Cache.Add(key, cacheEntry)
	c.mutex.Unlock()
	return nil
}

func (c *LocalCacheClient) PurgeCache() {
	c.Cache.Purge()
	fc.Info.Println("Purged proxy cache.")
}
