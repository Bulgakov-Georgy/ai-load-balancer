package cache

import (
	"context"
	"log"
	"time"

	"ai_load_balancer/internal/configuration"
	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

// InitRedis initializes the Redis client
func InitRedis() {
	config := configuration.Get()
	redisClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl,
		Password: config.RedisPassword,
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully!")
}

// GetString retrieves string value stored with the given key
func GetString(key string) (string, error) {
	value, err := redisClient.Get(ctx, key).Result()
	if err != nil && err.Error() != "redis: nil" {
		log.Println(err)
		return "", err
	}
	return value, nil
}

// GetInt retrieves int value stored with the given key
func GetInt(key string) (int, error) {
	value, err := redisClient.Get(ctx, key).Int()
	if err != nil && err.Error() != "redis: nil" {
		log.Println(err)
		return 0, err
	}
	return value, nil
}

// Set stores value in the cache by the given key
func Set(key string, value string, expiration time.Duration) {
	if err := redisClient.Set(ctx, key, value, expiration).Err(); err != nil {
		log.Println(err)
	}
}

// Increment increments value
func Increment(key string, expiration time.Duration) {
	redisClient.Incr(ctx, key)
	redisClient.Expire(ctx, key, expiration)
}
