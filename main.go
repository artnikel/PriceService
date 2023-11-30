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
	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"
)

func connectRedis() (*redis.Client, error) {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("could not parse config: ", err)
	}
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisPriceAddress,
		DB:   0,
	})
	_, err = client.Ping(client.Context()).Result()
	if err != nil {
		return nil, fmt.Errorf("error in method client.Ping(): %v", err)
	}
	return client, nil
}

// nolint gocritic
func main() {
	redisClient, err := connectRedis()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v", err)
	}
	defer func() {
		errClose := redisClient.Close()
		if errClose != nil {
			log.Fatalf("failed to disconnect from Redis: %v", errClose)
		}
	}()
	repoRedis := repository.NewRedisRepository(redisClient)
	servRedis := service.NewPriceService(repoRedis)
	handlRedis := handler.NewPriceHandler(servRedis)
	go servRedis.SubscribeAll(context.Background())
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("cannot create listener: %s", err)
	}
	fmt.Println("Price Service started")
	grpcServer := grpc.NewServer()
	proto.RegisterPriceServiceServer(grpcServer, handlRedis)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to serve listener: %s", err)
	}
}
