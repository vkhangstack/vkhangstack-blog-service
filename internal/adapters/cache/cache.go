package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0, // use default DB
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %v", err)
	}

	return &RedisCache{client: client}, nil
}

func (c *RedisCache) Get(key string, value interface{}) error {
	data, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return fmt.Errorf("cache miss for key %q", key)
	} else if err != nil {
		return fmt.Errorf("failed to get value for key %q: %v", key, err)
	}

	if err := json.Unmarshal([]byte(data), value); err != nil {
		return fmt.Errorf("failed to unmarshal cache value for key %q: %v", key, err)
	}

	return nil
}

func (c *RedisCache) Set(key string, value interface{}, duration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value for key %q: %v", key, err)
	}

	if err := c.client.Set(context.Background(), key, data, duration).Err(); err != nil {
		return fmt.Errorf("failed to set value for key %q: %v", key, err)
	}

	return nil
}

func (c *RedisCache) Delete(key string) error {
	if err := c.client.Del(context.Background(), key).Err(); err != nil {
		return fmt.Errorf("failed to delete value for key %q: %v", key, err)
	}
	return nil
}

func (c *RedisCache) SetString(key, value string, duration time.Duration) error {
	if err := c.client.Set(context.Background(), key, value, duration).Err(); err != nil {
		return fmt.Errorf("failed to set string value for key %q: %v", key, err)
	}
	return nil
}

func (c *RedisCache) GetString(key string) (string, error) {
	value, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("cache miss for key %q", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get string value for key %q: %v", key, err)
	}
	return value, nil
}

func (c *RedisCache) GetClient() *redis.Client {
	return c.client
}

func (c *RedisCache) Exists(key string) (bool, error) {
	exists, err := c.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %q: %v", key, err)
	}
	return exists > 0, nil
}

func (c *RedisCache) Ping() error {
	_, err := c.client.Ping(context.Background()).Result()
	return err
}

func (c *RedisCache) Close() error {
	return c.client.Close()
}
