// Package main of a project
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/artnikel/PriceService/internal/config"
	"github.com/artnikel/PriceService/internal/handler"
	"github.com/artnikel/PriceService/internal/repository"
	"github.com/artnikel/PriceService/internal/service"
	"github.com/artnikel/PriceService/proto"
	"github.com/caarlos0/env"
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

func connectRedis() (*redis.Client, error) {
	cfg := config.Variables{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisPriceAddress,
		Password: cfg.RedisPricePassword,
		DB:       0,
	})
	_, err := client.Ping(client.Context()).Result()
	if err != nil {
		return nil, fmt.Errorf("error in method client.Ping(): %v", err)
	}
	return client, nil
}

// nolint gocritic
func main() {
	redisClient, err := connectRedis()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer func() {
		errClose := redisClient.Close()
		if errClose != nil {
			log.Fatalf("Failed to disconnect from Redis: %v", errClose)
		}
	}()
	repoRedis := repository.NewRedisRepository(redisClient)
	servRedis := service.NewPriceService(repoRedis)
	handlRedis := handler.NewPriceHandler(servRedis)
	go servRedis.SendToAllSubscribedChans(context.Background())
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Cannot create listener: %s", err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterPriceServiceServer(grpcServer, handlRedis)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Failed to serve listener: %s", err)
	}
}
