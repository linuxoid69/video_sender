package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	KeyTTL      time.Duration
	RedisClient redis.Client
}

func NewClient(host, password string, db int, TTL time.Duration) *Client {
	return &Client{
		KeyTTL: TTL,
		RedisClient: *redis.NewClient(&redis.Options{
			Addr:     host,
			Password: password,
			DB:       db,
		}),
	}
}

func (c *Client) CreateJob(ctx context.Context, key string, value any) error {
	if err := c.Ping(ctx); err != nil {
		return err
	}

	if err := c.RedisClient.Set(ctx, key, value, c.KeyTTL).Err(); err != nil {
		return fmt.Errorf("can't set key: %w", err)
	}

	return nil
}

func (c *Client) GetJob(ctx context.Context, key string) (string, error) {
	if err := c.Ping(ctx); err != nil {
		return "", err
	}

	result, err := c.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("can't get key: %w", err)
	}

	return result, nil
}

func (c *Client) DeleteJob(ctx context.Context, keys ...string) error {
	if err := c.Ping(ctx); err != nil {
		return err
	}

	if err := c.RedisClient.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("can't delete key: %w", err)
	}

	return nil
}

func (c *Client) GetKeys(ctx context.Context, pattern string) ([]string, error) {
	if err := c.Ping(ctx); err != nil {
		return nil, err
	}

	keys := c.RedisClient.Keys(ctx, pattern)
	allKeys, err := keys.Result()
	if err != nil {
		return nil, fmt.Errorf("can't get keys: %w", err)
	}

	return allKeys, nil
}

func (c *Client) Ping(ctx context.Context) error {
	if _, err := c.RedisClient.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("can't to connect redis: %w", err)
	}

	return nil
}
