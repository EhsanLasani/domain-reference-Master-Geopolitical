// Package secrets implements secure secrets management with rotation
package secrets

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
)

type SecretsManager interface {
	GetSecret(ctx context.Context, key string) (string, error)
	RefreshSecrets(ctx context.Context) error
}

type EnvSecretsManager struct {
	cache map[string]cachedSecret
	mutex sync.RWMutex
}

type cachedSecret struct {
	value     string
	expiresAt time.Time
}

func NewEnvSecretsManager() SecretsManager {
	return &EnvSecretsManager{
		cache: make(map[string]cachedSecret),
	}
}

func (esm *EnvSecretsManager) GetSecret(ctx context.Context, key string) (string, error) {
	esm.mutex.RLock()
	if cached, exists := esm.cache[key]; exists && time.Now().Before(cached.expiresAt) {
		esm.mutex.RUnlock()
		return cached.value, nil
	}
	esm.mutex.RUnlock()
	
	return esm.fetchAndCache(ctx, key)
}

func (esm *EnvSecretsManager) fetchAndCache(ctx context.Context, key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("secret %s not found", key)
	}
	
	esm.mutex.Lock()
	esm.cache[key] = cachedSecret{
		value:     value,
		expiresAt: time.Now().Add(5 * time.Minute), // 5 minute cache
	}
	esm.mutex.Unlock()
	
	return value, nil
}

func (esm *EnvSecretsManager) RefreshSecrets(ctx context.Context) error {
	esm.mutex.Lock()
	defer esm.mutex.Unlock()
	
	// Clear cache to force refresh
	esm.cache = make(map[string]cachedSecret)
	return nil
}

// Helper functions for common secrets
func GetJWTSecret(ctx context.Context, sm SecretsManager) (string, error) {
	return sm.GetSecret(ctx, "JWT_SECRET")
}

func GetDatabasePassword(ctx context.Context, sm SecretsManager) (string, error) {
	return sm.GetSecret(ctx, "DB_PASSWORD")
}