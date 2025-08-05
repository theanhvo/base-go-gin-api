package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"baseApi/config"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

/* InitRedis initializes Redis connection */
func InitRedis(cfg *config.Config) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       0, // use default DB
	})

	// Test connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	log.Println("Redis connected successfully")
}

/* Set stores a value in Redis with expiration */
func Set(key string, value interface{}, expiration time.Duration) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return RedisClient.Set(ctx, key, json, expiration).Err()
}

/* Get retrieves a value from Redis */
func Get(key string, dest interface{}) error {
	val, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

/* Delete removes a key from Redis */
func Delete(key string) error {
	return RedisClient.Del(ctx, key).Err()
}

/* Exists checks if a key exists in Redis */
func Exists(key string) (bool, error) {
	count, err := RedisClient.Exists(ctx, key).Result()
	return count > 0, err
}

/* SetWithoutExpiration stores a value in Redis without expiration */
func SetWithoutExpiration(key string, value interface{}) error {
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return RedisClient.Set(ctx, key, json, 0).Err()
}

/* GetRedisClient returns the Redis client instance */
func GetRedisClient() *redis.Client {
	return RedisClient
}