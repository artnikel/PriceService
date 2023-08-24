// Package repository is a lower level of project
package repository

import (
	"context"
	"log"
	"strings"

	"github.com/artnikel/PriceService/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/shopspring/decimal"
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
func (r *RedisRepository) ReadPrices(ctx context.Context) (shares []*model.Share, e error) {
	tempMap := make(map[string]decimal.Decimal)
	for {
		result, err := r.client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{"shares", "0"},
			Count:   10,
			Block:   0,
		}).Result()
		if err != nil {
			log.Fatalf("RedisRepository-ReadPrices: Error when reading messages from Redis Stream:%v", err)
		}
		for _, message := range result {
			for _, msg := range message.Messages {
				data := msg.Values["message"].(string)
				parts := strings.Split(data, ":")
				if len(parts) != 2 {
					log.Fatalf("RedisRepository-ReadPrices: Incorrect message format: %s", data)
					continue
				}
				company := strings.TrimSpace(parts[0])
				priceStr := strings.TrimSpace(parts[1])

				price, err := decimal.NewFromString(priceStr)
				if err != nil {
					log.Fatalf("RedisRepository-ReadPrices: Error when converting price to number: %v", err)
					continue
				}
				tempMap[company] = price
				shares = append(shares, &model.Share{
					Company: company,
					Price:   price,
				})
			}
		}
		return shares, nil
	}

}
