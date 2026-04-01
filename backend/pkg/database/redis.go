package database

import (
	"context"
	"fmt"
	"log"
	"map-tile-system/pkg/config"

	"github.com/go-redis/redis/v8"
)

func InitRedis(cfg config.RedisConfig) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		return nil
	}

	log.Println("Redis connected successfully")
	return rdb
}
