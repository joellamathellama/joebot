package main

import (
	"fmt"

	"gopkg.in/redis.v4"
)

var (
	rc *redis.Client
)

// Connect to default port
func redisInit() {
	rc = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	fmt.Println("Redis Ping Pong test. Expect: 'PONG <nil>'")
	pong, err := rc.Ping().Result()
	fmt.Println(pong, err)
	if err != nil {
		writeErr(err)
		fmt.Println(err)
	}
	// Output: PONG <nil>

	// redist test: WORKING
	// redisSet(rc, "test key", "test string")
	// redisGet(rc, "test key")
}

func redisSet(c *redis.Client, key string, value string) {
	err := c.Set(key, value, 0).Err()
	if err != nil {
		writeErr(err)
		fmt.Println(err)
	}
}

func redisGet(c *redis.Client, key string) (val string) {
	val, err := c.Get(key).Result()
	if err != nil {
		writeErr(err)
		fmt.Println(err)
	}
	return
}
