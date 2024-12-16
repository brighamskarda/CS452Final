package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func CheckRedisCache(fromArticle string, toArticle string) string {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // No password set
		DB:       0,  // Use default DB
		Protocol: 2,  // Connection protocol
	})

	ctx := context.Background()

	val, err := client.Get(ctx, fromArticle+"->"+toArticle).Result()
	if err != nil {
		return ""
	}
	return val
}
