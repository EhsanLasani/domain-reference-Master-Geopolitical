// ============================================================================
// FILE: repository_catalogue.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Data Access Layer - Catalogue
// PURPOSE: Provide business layer with available operations catalogue
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package catalogue

import (
	"context"
	"time"
)

// RepositoryCatalogue provides business layer with available operations
type RepositoryCatalogue struct {
	Countries CountryOperations
	Regions   RegionOperations
	Languages LanguageOperations
}

// CountryOperations defines all available country operations
type CountryOperations struct {
	// CRUD Operations
	Create      func(ctx context.Context, country interface{}) error
	GetByID     func(ctx context.Context, id string) (interface{}, error)
	GetByCode   func(ctx context.Context, code string) (interface{}, error)
	Update      func(ctx context.Context, country interface{}) error
	Delete      func(ctx context.Context, id string) error
	
	// Query Operations
	ListActive  func(ctx context.Context) ([]interface{}, error)
	Search      func(ctx context.Context, query string) ([]interface{}, error)
	
	// Business Operations
	ExistsByCode    func(ctx context.Context, code string) (bool, error)
	GetActiveCount  func(ctx context.Context) (int, error)
	
	// Validation Operations
	ValidateCode    func(code string) error
	ValidateName    func(name string) error
}

// RegionOperations defines all available region operations
type RegionOperations struct {
	Create         func(ctx context.Context, region interface{}) error
	GetByCode      func(ctx context.Context, code string) (interface{}, error)
	GetChildren    func(ctx context.Context, parentID string) ([]interface{}, error)
}

// LanguageOperations defines all available language operations  
type LanguageOperations struct {
	Create           func(ctx context.Context, language interface{}) error
	GetByCode        func(ctx context.Context, code string) (interface{}, error)
	GetByDirection   func(ctx context.Context, direction string) ([]interface{}, error)
}

// OperationMetadata provides information about each operation
type OperationMetadata struct {
	Name            string        `json:"name"`
	Description     string        `json:"description"`
	Parameters      []Parameter   `json:"parameters"`
	ReturnType      string        `json:"return_type"`
	ErrorCodes      []string      `json:"error_codes"`
	CacheEnabled    bool          `json:"cache_enabled"`
	CacheTTL        time.Duration `json:"cache_ttl"`
}

type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

// GetOperationCatalogue returns metadata for all available operations
func GetOperationCatalogue() map[string]OperationMetadata {
	return map[string]OperationMetadata{
		"Countries.GetByCode": {
			Name:        "Get Country by Code",
			Description: "Retrieve country using ISO 3166-1 alpha-2 code",
			Parameters: []Parameter{
				{Name: "code", Type: "string", Required: true, Description: "ISO country code"},
			},
			ReturnType:   "Country",
			ErrorCodes:   []string{"GEO-1001", "GEO-1003"},
			CacheEnabled: true,
			CacheTTL:     24 * time.Hour,
		},
		"Countries.ListActive": {
			Name:        "List Active Countries",
			Description: "Retrieve all active countries",
			Parameters:  []Parameter{},
			ReturnType:  "[]Country",
			ErrorCodes:  []string{"GEO-9001"},
			CacheEnabled: true,
			CacheTTL:     6 * time.Hour,
		},
	}
}