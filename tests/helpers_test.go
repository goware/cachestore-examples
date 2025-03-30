package cachestore_e2e_test

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func redisFlushAll() {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 9})
	redisClient.FlushAll(context.Background())
}
