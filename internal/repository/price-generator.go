package repository

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

func GeneratePrice(client *redis.Client, timeSecondSleep time.Duration) {
	ctx := context.Background()
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))

	prices := map[string]int{
		"Logitech": 10,
		"Apple": 90,
		"Microsoft": 150,
	}
	for {
		for company, price := range prices {
			change := rng.Intn(7) - 3
			price += change
			_, err := client.XAdd(ctx, &redis.XAddArgs{
				Stream: "messagestream",
				Values: map[string]interface{}{
					"message": fmt.Sprintf("%s: %d", company, price),
				},
			}).Result()
			if err != nil {
				log.Fatalf("Error when writing a message to Redis Stream: %v", err)
			}
			prices[company] = price
		}
		time.Sleep(timeSecondSleep * time.Second)
	}
}