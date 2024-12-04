//go:generate mockgen -typed -source=redis.go -destination=redis_mock.go -package=database
package database

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type RedisClientInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration int) error
	Delete(ctx context.Context, key string) error
	Close() error
}

type RedisClient struct {
	client redis.UniversalClient
}

func (c *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration int) error {
	return c.client.Set(ctx, key, value, time.Duration(expiration)).Err()
}

func (c *RedisClient) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

func (c *RedisClient) Close() error {
	return c.client.Close()
}

func NewRedisClient(addr, password string) *RedisClient {
	client := &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
		}),
	}

	// check connection
	pong, err := client.client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("cailed to connect to Redis: %v", err)
	}

	fmt.Println("connected to Redis:", pong)
	return client
}
