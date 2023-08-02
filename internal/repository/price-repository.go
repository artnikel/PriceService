// Package repository is a lower level of project
package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/artnikel/PriceService/internal/model"
	"github.com/go-redis/redis/v8"
)

// RedisRepository contains objects of type *redis.Client
type RedisRepository struct {
	client *redis.Client
}

// NewRedisRepository accepts an object of *redis.Client and returns an object of type *Redis
func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
	}
}

// nolint gonmd
// ReadPrices is a method that read price by field company from redis stream adn returns price of this company
func (r *RedisRepository) ReadPrices(ctx context.Context) (actions []*model.Action, e error) {
	result, err := r.client.XRevRange(ctx, "messagestream", "+", "-").Result()
	if err != nil {
		return nil, fmt.Errorf("PriceServiceRepository-ReadFromStream-XRevRange: error: %w", err)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("PriceServiceRepository-ReadFromStream: error: message is empty")
	}
	err = json.Unmarshal([]byte(result[0].Values["message"].(string)), &actions)
	if err != nil {
		return nil, fmt.Errorf("PriceServiceRepository-ReadFromStream: error in method json.Unmarshal: %w", err)
	}
	return actions, nil
}
