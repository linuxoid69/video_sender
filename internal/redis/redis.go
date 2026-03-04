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
		RedisClient: *redis.NewClient(
			&redis.Options{
				Addr:     host,
				Password: password,
				DB:       db,
			}),
	}
}

func (c *Client) Create(ctx context.Context, key string, value any) (err error) {
	if err = c.Ping(ctx); err != nil {
		return err
	}

	if err = c.RedisClient.Set(ctx, key, value, c.KeyTTL).Err(); err != nil {
		return fmt.Errorf("failed to set key %s with value %v: %w", key, value, err)
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key string) (result string, err error) {
	if err = c.Ping(ctx); err != nil {
		return "", err
	}

	result, err = c.RedisClient.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}

	return result, nil
}

func (c *Client) Delete(ctx context.Context, keys ...string) (err error) {
	if err = c.Ping(ctx); err != nil {
		return err
	}

	if err = c.RedisClient.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", keys, err)
	}

	return nil
}

func (c *Client) Keys(ctx context.Context, pattern string) (allKeys []string, err error) {
	if err = c.Ping(ctx); err != nil {
		return nil, err
	}

	allKeys, err = c.RedisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys by pattern %s: %w", pattern, err)
	}

	return allKeys, nil
}

func (c *Client) Ping(ctx context.Context) (err error) {
	if _, err = c.RedisClient.Ping(ctx).Result(); err != nil {
		return fmt.Errorf("failed to connect redis: %w", err)
	}

	return nil
}
