package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var rdb *redis.Client

func NewRedisClient(addr string, password string, db int) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis-dev:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func Set(value string) error {

	if rdb == nil {
		NewRedisClient("localhost:6379", "", 0)
	}

	res1, err := rdb.RPush(ctx, "limits:log", "limit:"+value).Result()

	if err != nil {
		log.Println("#######################333", err)
		return err
	}
	fmt.Println("@@@@@@@@@@@@@@@@@@@@", res1) // >>> 1
	return nil
}

// func Get(key string) (string, error) {
// 	val, err := rdb.Get(ctx, key).Result()
// 	if err != nil {
// 		return "", err
// 	}
// 	fmt.Println("key", val)
// 	return val, nil
// }

// func ExampleClient() {

// 	val2, err := rdb.Get(ctx, "key2").Result()
// 	if err == redis.Nil {
// 		fmt.Println("key2 does not exist")
// 	} else if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Println("key2", val2)
// 	}
// 	// Output: key value
// 	// key2 does not exist
// }
