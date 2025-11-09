package flags

import (
	"context"
	"sync"
)

// FeatureFlags implements guideline 11 - Feature Flags
type FeatureFlags struct {
	flags map[string]bool
	mu    sync.RWMutex
}

func NewFeatureFlags() *FeatureFlags {
	return &FeatureFlags{
		flags: make(map[string]bool),
	}
}

func (f *FeatureFlags) IsEnabled(ctx context.Context, flag string, tenantID string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	// Check tenant-specific flag first
	tenantFlag := flag + ":" + tenantID
	if enabled, exists := f.flags[tenantFlag]; exists {
		return enabled
	}
	
	// Fall back to global flag
	return f.flags[flag]
}

func (f *FeatureFlags) SetFlag(flag string, enabled bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.flags[flag] = enabled
}

func (f *FeatureFlags) SetTenantFlag(flag, tenantID string, enabled bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.flags[flag+":"+tenantID] = enabled
}