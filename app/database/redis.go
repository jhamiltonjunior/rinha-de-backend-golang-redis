package database

import "github.com/redis/go-redis/v9"

func InitializeRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "redis_cache:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})
}