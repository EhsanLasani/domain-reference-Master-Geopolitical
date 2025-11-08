// ============================================================================
// FILE: country_model.go
// PACKAGE: models
// DOMAIN: Reference Master Geopolitical
// LAYER: Data Access Layer - ORM/ODM Abstractions
// VERSION: 1.0.0
// CREATED: 2025-11-07
// LAST_MODIFIED: 2025-11-07
// ============================================================================
//
// PURPOSE:
//   Domain model for Country entity representing sovereign nations and territories
//   with complete geopolitical information following international standards.
//
// BUSINESS CONTEXT:
//   Core reference entity for all country-related operations across LASANI platform.
//   Serves as foundational data for user management, tenant management, and commerce domains.
//
// COMPLIANCE:
//   - ISO 3166-1 standards for country codes and identification
//   - LASANI audit system integration (when extended)
//   - International collation support for proper text sorting
//
// PERFORMANCE CHARACTERISTICS:
//   - Read-heavy workload (95% reads, 5% writes)
//   - Average entity size: 1KB in memory
//   - Expected query volume: 1000+ requests/second
//   - Cache-friendly (24-hour TTL recommended)
//
// INTEGRATION POINTS:
//   - Referenced by: UserProfiles, TenantSettings, BillingAddresses
//   - References: Regions (optional), Languages (optional)
//   - Events: CountryCreated, CountryUpdated, CountryDeactivated
//
// AUTHOR: Development Team
// REVIEWER: Database Architect
// ============================================================================

package models

import (
	"time"
)

// Country represents a sovereign nation or territory with complete geopolitical information.
//
// This entity serves as the foundational reference for all country-related data across
// the LASANI platform. It follows ISO 3166-1 standards for country identification
// and supports international business operations.
//
// Business Rules:
//   - CountryCode must be valid ISO 3166-1 alpha-2 format (e.g., "US", "GB")
//   - CountryName is required and must be 2-100 characters
//   - ISO3Code must be valid ISO 3166-1 alpha-3 format (e.g., "USA", "GBR")
//   - OfficialName provides the formal government name
//   - IsActive controls business operation visibility
//
// Performance Notes:
//   - Optimized for read operations with minimal memory footprint
//   - Suitable for caching with 24-hour TTL
//   - Indexed by CountryCode for O(1) lookup performance
//
// Example Usage:
//
//	country := &Country{
//		CountryCode:  "US",
//		CountryName:  "United States",
//		ISO3Code:     "USA",
//		OfficialName: "United States of America",
//		IsActive:     true,
//	}
//
type Country struct {
	// Primary Key - UUID for unique identification
	CountryID string `json:"country_id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	
	// Business identity - ISO 3166-1 alpha-2 country code (required)
	CountryCode string `json:"country_code" gorm:"uniqueIndex:idx_country_code_tenant;size:2" validate:"required,len=2,alpha"`
	
	// Display name - human-readable country name (required)
	CountryName string `json:"country_name" gorm:"size:100" validate:"required,min=2,max=100"`
	
	// Extended identification - ISO 3166-1 alpha-3 code (optional)
	ISO3Code string `json:"iso3_code" gorm:"size:3" validate:"omitempty,len=3,alpha"`
	
	// Formal name - official government designation (optional)
	OfficialName string `json:"official_name" gorm:"size:200" validate:"omitempty,max=200"`
	
	// Status management - controls business operation visibility
	IsActive bool `json:"is_active" gorm:"default:true"`
	IsDeleted bool `json:"is_deleted" gorm:"default:false"`
	
	// LASANI Audit Fields - Multi-tenant compliance
	TenantID     string    `json:"tenant_id" gorm:"uniqueIndex:idx_country_code_tenant;size:36" validate:"required"`
	CreatedBy    string    `json:"created_by" gorm:"size:36"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	ModifiedBy   string    `json:"modified_by" gorm:"size:36"`
	ModifiedAt   time.Time `json:"modified_at" gorm:"autoUpdateTime"`
	Version      int       `json:"version" gorm:"default:1"`
	ChangeReason string    `json:"change_reason" gorm:"size:500"`
}