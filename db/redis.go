package db

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis_rate/v10"
	limiter "github.com/go-redis/redis_rate/v10"
	redis "github.com/redis/go-redis/v9"
	"kgoel085.com/url-shortner/config"
)

var (
	RedisClient  *redis.Client
	RedisLimiter *redis_rate.Limiter
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

	RedisLimiter = redis_rate.NewLimiter(RedisClient)

}

func CheckRateLimitInTimeUnit(ctx context.Context, key string, requests int, timeUnit time.Duration) (bool, error) {
	var limit limiter.Limit
	switch timeUnit {
	case time.Minute:
		limit = limiter.PerMinute(requests)
	case time.Hour:
		limit = limiter.PerHour(requests)
	case time.Second:
		limit = limiter.PerSecond(requests)
	default:
		return false, fmt.Errorf("unsupported time unit")
	}
	return checkRateLimiter(ctx, key, limit)
}

func checkRateLimiter(ctx context.Context, key string, limit limiter.Limit) (bool, error) {
	res, err := RedisLimiter.Allow(ctx, key, limit)
	if err != nil {
		return false, err
	}
	return res.Allowed > 0, nil
}
