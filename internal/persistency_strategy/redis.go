package persistencystrategy

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

var rdb *redis.Client

func (l PersistencyStrategyStruct) Init(config []string) {
	NewRedisClient(config)
}

func (l PersistencyStrategyStruct) Log(msg string) error {
	err := RPush(msg)
	if err != nil {
		log.Println("Error pushing to Redis:", err)
		return err
	}
	log.Println("Message pushed to Redis:", msg)
	return nil
}

func NewRedisClient(config []string) {

	rdb = redis.NewClient(&redis.Options{})

	rdb = redis.NewClient(&redis.Options{
		Addr:     config[0] + ":6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	rdb.Ping(ctx)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Println("Error connecting to Redis:", err)
	} else {
		log.Println("Connected to Redis")
	}

}

func RPush(value string) error {

	if rdb == nil {
		log.Println("Redis client not initialized")
		return fmt.Errorf("redis client not initialized")
	}

	res1, err := rdb.RPush(ctx, "limits:log", "limit:"+value).Result()

	if err != nil {
		log.Println("#######################", err)
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
