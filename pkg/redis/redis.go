package redis

import (
	"context"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
)

func InitializeRedis() *redis.Client {
	conn := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_CONNECTION_STRING"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// checking if redis is connected
	pong, err := conn.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Redis Connection Failed",
			err)
	}

	log.Println("Redis Successfully Connected.",
		"Ping", pong)

	return conn
}
