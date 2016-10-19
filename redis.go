package main

import (
	"fmt"

	"gopkg.in/redis.v4"
)

var (
	redisClient *redis.Client
)

// Connect to default port
func redisInit() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := redisClient.Ping().Result()
	fmt.Println(pong, err)
	// Output: PONG <nil>

	// redist test: WORKING
	// redisSet(redisClient, "test key", "test string")
	// redisGet(redisClient, "test key")
}

func redisSet(c *redis.Client, key string, value string) {
	err := c.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}

func redisGet(c *redis.Client, key string) {
	val, err := c.Get(key).Result()
	if err != nil {
		// panic(err)
		fmt.Println("Invalid Key")
	}
	fmt.Println(key, val)
}
