package rds

import (
	"fmt"
	"gopkg.in/redis.v4"
	"joebot/tools"
)

var (
	RC *redis.Client
)

// Connect to default port
func RedisInit() {
	RC = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	fmt.Println("Redis Ping Pong test. Expect: 'PONG <nil>'")
	pong, err := RC.Ping().Result()
	fmt.Println(pong, err)
	if err != nil {
		tools.WriteErr(err)
		fmt.Println(err)
	}
	// Output: PONG <nil>

	// redist test: WORKING
	// rds.RedisSet(RC, "test key", "test string")
	// rds.RedisGet(RC, "test key")
}

func RedisSet(c *redis.Client, key string, value string) bool {
	err := c.Set(key, value, 0).Err()
	if err != nil {
		tools.WriteErr(err)
		return false
	}
	return true
}

func RedisGet(c *redis.Client, key string) (val string, err error) {
	val, err = c.Get(key).Result()
	if err != nil {
		tools.WriteErr(err)
	}
	return
}

func RedisLPush(c *redis.Client, key string, value string) bool {
	err := c.LPush(key, value).Err()
	if err != nil {
		tools.WriteErr(err)
		return false
	}
	return true
}

func RedisLRem(c *redis.Client, key string, value string) bool {
	err := c.LRem(key, 0, value).Err()
	if err != nil {
		tools.WriteErr(err)
		return false
	}
	return true
}

func RedisLRange(c *redis.Client, key string) (val []string) {
	val, err := c.LRange(key, 0, -1).Result()
	if err != nil {
		tools.WriteErr(err)
	}
	return
}
