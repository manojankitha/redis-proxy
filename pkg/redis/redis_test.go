package redis

import (
	"log"
	"testing"
	"time"

	"github.com/freeconf/yang/fc"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
)

// Test to confirm redis client is able to set key value
func TestRedisClientSetKey(t *testing.T) {
	redisClient, err := RedisInit("localhost:6379")
	if err != nil {
		log.Fatalf("Error initializing redis service: %v", err.Error())
	}
	fc.Info.Println(t.Name())
	key1 := "testKey"
	value1 := "testValue"
	err = redisClient.SetKey(key1, value1, time.Minute*10)
	if err != nil {
		log.Fatalf("Error setting key: %v", err.Error())
	}
	redisClient.FlushRedisDB()
}

// Test to confirm redis client is able to get value of key
func TestRedisClientGetKey(t *testing.T) {

	redisClient, err := RedisInit("localhost:6379")
	if err != nil {
		log.Fatalf("Error initializing redis service: %v", err.Error())
	}
	fc.Info.Println(t.Name())
	key1 := "testKey"
	value1 := "testValue"
	key2 := "faultkey"

	// check that new key that is not set returns error
	_, err = redisClient.GetKey(key2)
	assert.Error(t, err)
	assert.Equal(t, err, redis.Nil)

	err = redisClient.SetKey(key1, value1, time.Minute*10)
	if err != nil {
		log.Fatalf("Error setting key: %v", err.Error())
	}
	// check value of stored key
	val, err := redisClient.GetKey(key1)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}
	assert.Equal(t, value1, val)
	redisClient.FlushRedisDB()
}
