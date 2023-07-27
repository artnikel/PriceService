package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

func (r *RedisRepository) ReadPrices(ctx context.Context, company string) (float64, error) {
	var price float64

	result, err := r.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"messagestream", "0"},
		Count:   10,
		Block:   0,
	}).Result()
	if err != nil {
		return 0, fmt.Errorf("error when reading messages from Redis Stream: %w", err)
	}

	for _, message := range result {
		for _, msg := range message.Messages {
			data := msg.Values["message"].(string)
			parts := strings.Split(data, ":")
			if len(parts) != 2 {
				return 0, fmt.Errorf("incorrect message format: %s", data)
			}

			name := strings.TrimSpace(parts[0])
			priceStr := strings.TrimSpace(parts[1])

			if name == company {
				price, err = strconv.ParseFloat(priceStr, 64)
				if err != nil {
					return 0, fmt.Errorf("error when converting price to number: %w", err)
				}
				return price, nil
			}
		}
	}

	return 0, fmt.Errorf("company not found")
}
