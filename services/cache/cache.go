package cache

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() error {
	redisURL := os.Getenv("REDIS_URL")
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}
	RedisClient = redis.NewClient(opt)
	return nil
}

func GetCache(ctx context.Context, redisClient *redis.Client, key string, dest any) (bool, error) {
	data, err := redisClient.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil // Cache miss
	}
	if err != nil {
		return false, err // Some other error
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, err
	}
	return true, nil
}

func DeleteCache(ctx context.Context, redisClient *redis.Client, key string) error {
	return redisClient.Del(ctx, key).Err()
}

func SetCache(ctx context.Context, redisClient *redis.Client, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, key, data, 5*time.Minute).Err()
}
