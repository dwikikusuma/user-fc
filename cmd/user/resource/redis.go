package resource

import (
	"commerce/config"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
)

var RedisClient *redis.Client

func InitRedis(cfg *config.Config) *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
	})

	// Test the connection
	ctx := context.Background()
	pingResponse, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Redis connection successful:", pingResponse)

	return RedisClient
}
