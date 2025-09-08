package db

import (
	"context"
	"fmt"

	"example.com/url-shortner/config"
	redis "github.com/go-redis/redis/v9"
)

var (
	RedisClient *redis.Client
)

func InitRedis() {
	addr := config.Config.REDIS.Addr
	password := config.Config.REDIS.Password
	db := config.Config.REDIS.DB

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	redisPing, redisErr := RedisClient.Ping(context.Background()).Result()
	if redisErr != nil {
		panic(redisErr)
	}
	fmt.Println("Redis Ping:", redisPing)
}
