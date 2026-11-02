package ledger

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(ctx context.Context) error {
	dbNumber := 0
	if value := os.Getenv("REDIS_DB"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid REDIS_DB: %w", err)
		}
		dbNumber = parsed
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     envOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       dbNumber,
	})

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("ping redis: %w", err)
	}
	return nil
}

func CloseRedis() error {
	if redisClient == nil {
		return nil
	}
	return redisClient.Close()
}
