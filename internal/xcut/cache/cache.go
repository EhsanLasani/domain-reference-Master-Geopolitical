// Package cache implements tenant-scoped caching framework
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Invalidate(ctx context.Context, pattern string) error
}

type TenantCache struct {
	cache    Cache
	tenantID string
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) *RedisCache {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	
	return &RedisCache{client: rdb}
}

func (r *RedisCache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("key not found")
	}
	if err != nil {
		return nil, fmt.Errorf("cache get error: %w", err)
	}
	return []byte(val), nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) Invalidate(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}
	return nil
}

func NewTenantCache(cache Cache, tenantID string) *TenantCache {
	return &TenantCache{
		cache:    cache,
		tenantID: tenantID,
	}
}

func (tc *TenantCache) scopedKey(key string) string {
	return fmt.Sprintf("tenant:%s:%s", tc.tenantID, key)
}

func (tc *TenantCache) Get(ctx context.Context, key string) ([]byte, error) {
	return tc.cache.Get(ctx, tc.scopedKey(key))
}

func (tc *TenantCache) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return tc.cache.Set(ctx, tc.scopedKey(key), value, ttl)
}

func (tc *TenantCache) Delete(ctx context.Context, key string) error {
	return tc.cache.Delete(ctx, tc.scopedKey(key))
}

func (tc *TenantCache) InvalidateTenant(ctx context.Context) error {
	pattern := fmt.Sprintf("tenant:%s:*", tc.tenantID)
	return tc.cache.Invalidate(ctx, pattern)
}