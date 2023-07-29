package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/require"
)

var rpcRedis *RedisRepository

func SetupTestRedis() (*redis.Client, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}
	resource, err := pool.Run("redis", "latest", []string{})
	if err != nil {
		return nil, nil, fmt.Errorf("could not start resource: %w", err)
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		DB:   0,
	})
	ctx := context.Background()
	err = pool.Retry(func() error {
		var pong string
		pong, err = redisClient.Ping(ctx).Result()
		if err != nil {
			return fmt.Errorf("error in method redisClient.Ping(): %w", err)
		}
		if pong != "PONG" {
			return fmt.Errorf("unexpected response from Redis: %s", pong)
		}
		return nil
	})
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to Redis: %w", err)
	}
	cleanup := func() {
		redisClient.FlushDB(ctx)
		redisClient.Close()
		pool.Purge(resource)
	}
	return redisClient, cleanup, nil
}

func TestMain(m *testing.M) {
	rdsClient, cleanupRds, err := SetupTestRedis()
	if err != nil {
		fmt.Println(err)
		cleanupRds()
		os.Exit(1)
	}
	rpcRedis = NewRedisRepository(rdsClient)

	exitCode := m.Run()

	cleanupRds()
	os.Exit(exitCode)
}

func TestReadNotExist(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := rpcRedis.ReadPrices(ctx, "Apple")
	require.Error(t, err)
}

func TestReadPrices(t *testing.T) {
	var (
		testCompany = "Amazon"
		testPrice   = 180.74
	)
	_, err := rpcRedis.client.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "messagestream",
		Values: map[string]interface{}{
			"message": fmt.Sprintf("%s: %.2f", testCompany, testPrice),
		},
		MaxLen: 5,
	}).Result()
	require.NoError(t, err)
	respPrice, err := rpcRedis.ReadPrices(context.Background(), testCompany)
	require.NoError(t, err)
	require.Equal(t, respPrice, testPrice)
}
