package repository

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

var prices sync.Map

type RedisRepository struct {
	client    *redis.Client
	stockData *sync.Map
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client:    client,
		stockData: &sync.Map{},
	}
}

func (r *RedisRepository) GeneratePrices(ctx context.Context, client *redis.Client) error {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))

	initialPrices := map[string]float64{
		"Logitech":  172.3,
		"Apple":     930.6,
		"Microsoft": 859.5,
		"Samsung":   565.3,
		"Xerox":     415.7,
	}

	for company, price := range initialPrices {
		r.stockData.Store(company, price)
	}

	for {
		r.stockData.Range(func(key, value interface{}) bool {
			company := key.(string)
			price := value.(float64)

			change := rng.Float64()*40.0 - 20.0
			price += change

			if price < 0 {
				price = 0.1
			}

			r.stockData.Store(company, price)

			_, err := r.client.XAdd(ctx, &redis.XAddArgs{
				Stream: "messagestream",
				Values: map[string]interface{}{
					"message": fmt.Sprintf("%s: %.2f", company, price),
				},
				MaxLen: 5,
			}).Result()
			if err != nil {
				log.Fatalf("Error when writing a message to Redis Stream: %v", err)
			}

			return true
		})
		time.Sleep(time.Second / 2)
	}
}

func (r *RedisRepository) ReadPrices(ctx context.Context, client *redis.Client) (map[string]float64, error) {
	tempMap := make(map[string]float64)

	for {
		result, err := r.client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{"messagestream", "0"},
			Count:   10,
			Block:   0,
		}).Result()

		if err != nil {
			log.Fatalf("Error when reading messages from Redis Stream:%v", err)
		}

		for _, message := range result {
			for _, msg := range message.Messages {
				data := msg.Values["message"].(string)
				parts := strings.Split(data, ":")
				if len(parts) != 2 {
					log.Fatalf("Incorrect message format: %s", data)
					continue
				}

				company := strings.TrimSpace(parts[0])
				priceStr := strings.TrimSpace(parts[1])

				price, err := strconv.ParseFloat(priceStr, 64)
				if err != nil {
					log.Fatalf("Error when converting price to number: %v", err)
					continue
				}
				tempMap[company] = price

				_, err = r.client.XDel(ctx, "messagestream", msg.ID).Result()
				if err != nil {
					log.Fatalf("Error when deleting message from Redis Stream: %v", err)
				}
			}
		}
		for company, price := range tempMap {
			prices.Store(company, price)
		}
		resultMap := make(map[string]float64)
		for k, v := range tempMap {
			resultMap[k] = v
		}
		return resultMap, nil
	}
}
