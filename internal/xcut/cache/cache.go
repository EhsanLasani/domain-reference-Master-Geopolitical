package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache interface defines caching operations
type Cache interface {
	Set(ctx context.Context, tenantID, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, tenantID, key string, dest interface{}) error
	Delete(ctx context.Context, tenantID, key string) error
}

// RedisCache implements Cache interface
type RedisCache struct {
	client *redis.Client
}

func NewCache(addr, password string, db int) Cache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{client: rdb}
}

func (c *RedisCache) Set(ctx context.Context, tenantID, key string, value interface{}, ttl time.Duration) error {
	scopedKey := c.scopeKey(tenantID, key)
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(ctx, scopedKey, data, ttl).Err()
}

func (c *RedisCache) Get(ctx context.Context, tenantID, key string, dest interface{}) error {
	scopedKey := c.scopeKey(tenantID, key)
	data, err := c.client.Get(ctx, scopedKey).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

func (c *RedisCache) Delete(ctx context.Context, tenantID, key string) error {
	scopedKey := c.scopeKey(tenantID, key)
	return c.client.Del(ctx, scopedKey).Err()
}

func (c *RedisCache) scopeKey(tenantID, key string) string {
	return fmt.Sprintf("tenant:%s:geo:%s", tenantID, key)
}