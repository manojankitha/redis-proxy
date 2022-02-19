package redis

import (
	"log"
	"time"

	"github.com/freeconf/yang/fc"
	"github.com/go-redis/redis"
)

var (
	client = &RedisClient{}
)

type RedisClient struct {
	C *redis.Client
}

//GetClient get the redis client
func RedisInit(redisAddr string) (*RedisClient, error) {
	c := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := c.Ping().Err(); err != nil {
		log.Printf("Unable to connect to redis " + err.Error())
		return nil, err
	}
	fc.Info.Println("Initializing redis client...")
	client.C = c
	return client, nil
}

//GetKey get key
func (client *RedisClient) GetKey(key string) (string, error) {
	val, err := client.C.Get(key).Result()
	if err == redis.Nil || err != nil {
		log.Printf("Redis was unable to retrive key. Error details: %v.", err)
		return "", err
	}

	if err != nil {
		return "", err
	}
	return val, nil
}

//SetKey set key
func (client *RedisClient) SetKey(key string, value string, expiration time.Duration) error {
	err := client.C.Set(key, value, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (client *RedisClient) FlushRedisDB() {
	client.C.FlushDB()
	fc.Info.Println("Flushed Redis DB.")
}
