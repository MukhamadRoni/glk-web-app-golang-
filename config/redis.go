package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	redisHost := GetEnv("REDIS_HOST", "localhost:6379")
	redisPassword := GetEnv("REDIS_PASSWORD", "")
	
	RDB = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword, // no password set
		DB:       0,             // use default DB
	})

	_, err := RDB.Ping(Ctx).Result()
	if err != nil {
		log.Println("[Redis] Failed to connect:", err)
	} else {
		log.Println("[Redis] Connected to", redisHost)
	}
}
