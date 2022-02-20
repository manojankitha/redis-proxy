package cache

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/freeconf/yang/fc"
	"github.com/stretchr/testify/assert"
)

const CACHE_CAPACITY = 5
const CACHE_EXPIRY_TIME = 500 // in ms

// Test to check if initialization of proxy cache works
func TestCacheInitialize(t *testing.T) {
	c, err := CacheInit(CACHE_CAPACITY, CACHE_EXPIRY_TIME)
	if err != nil {
		log.Fatalf("Error initializing proxy cache: %v", err.Error())
	}
	fc.Info.Println(t.Name())
	assert.Equal(t, c.cacheExpiryTime, int64(CACHE_EXPIRY_TIME))

	c.PurgeCache()

}

// Test to check hasExpired function
func TestHasExpired(t *testing.T) {
	c, err := CacheInit(CACHE_CAPACITY, CACHE_EXPIRY_TIME)
	if err != nil {
		log.Fatalf("Error initializing proxy cache: %v", err.Error())
	}
	fc.Info.Println(t.Name())

	// test if value persists before expiration time
	c.SetKey("testPersistkey", "testPersistValue")
	time.Sleep(480 * time.Millisecond)
	value, ok := c.GetKey("testPersistkey")
	assert.Equal(t, true, ok)
	assert.Equal(t, "testPersistValue", value)

	// Test if value got evicted from cache after expiration
	c.SetKey("testExpiredkey", "testExpiredValue")
	time.Sleep(501 * time.Millisecond)
	value, ok = c.GetKey("testExpiredkey")
	assert.Equal(t, false, ok)
	assert.Equal(t, "", value)

	c.PurgeCache()

}

// Test to check addition of new element to cache
func TestAddElementToCache(t *testing.T) {
	c, err := CacheInit(CACHE_CAPACITY, CACHE_EXPIRY_TIME)
	if err != nil {
		log.Fatalf("Error initializing proxy cache: %v", err.Error())
	}
	fc.Info.Println(t.Name())
	c.SetKey("testKey", "testValue")
	c.SetKey("testKey", "testValueChanged")
	value, ok := c.GetKey("testKey")
	assert.Equal(t, true, ok)
	// checks if new value got stored for same key
	assert.Equal(t, "testValueChanged", value)

	c.PurgeCache()
}

// Test to check getKey function handles keys that are not set in the first place graciously and return ""
func TestGetKeyCornerCase(t *testing.T) {
	c, err := CacheInit(CACHE_CAPACITY, CACHE_EXPIRY_TIME)
	if err != nil {
		log.Fatalf("Error initializing proxy cache: %v", err.Error())
	}
	fc.Info.Println(t.Name())
	value, ok := c.GetKey("ImaginaryKey")
	assert.Equal(t, false, ok)
	assert.Equal(t, "", value)

	c.PurgeCache()

}

// Test to check cache Eviction policy - LRU
func TestCacheEviction(t *testing.T) {
	c, err := CacheInit(CACHE_CAPACITY, CACHE_EXPIRY_TIME)
	if err != nil {
		log.Fatalf("Error initializing proxy cache: %v", err.Error())
	}
	fc.Info.Println(t.Name())

	c.SetKey("DummyKey", "DummyValue")
	for i := 0; i < CACHE_CAPACITY; i++ {
		c.SetKey(strconv.Itoa(i), strconv.Itoa(i))

	}
	// get Key [0] once
	c.GetKey(strconv.Itoa(0))

	// get keys [1:CACHE_CAPACITY] twice
	for i := 1; i < CACHE_CAPACITY; i++ {
		c.GetKey(strconv.Itoa(i))
		c.GetKey(strconv.Itoa(i))
	}

	// at this point the DummyKey should be evicted because it's recently used frequency is 0
	value, ok := c.GetKey("DummyKey")
	assert.Equal(t, false, ok)
	assert.Equal(t, "", value)
	c.PurgeCache()
}
