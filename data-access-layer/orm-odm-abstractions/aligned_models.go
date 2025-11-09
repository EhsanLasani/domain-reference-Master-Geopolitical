package models

import (
	"time"
	"net"
	"encoding/json"
	"github.com/google/uuid"
)

// Country represents the country entity with full LASANI audit compliance
type Country struct {
	// Primary Key
	CountryID         uuid.UUID `json:"country_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	
	// Business Fields
	CountryCode       string     `json:"country_code" gorm:"type:char(2);uniqueIndex;not null"`
	CountryName       string     `json:"country_name" gorm:"type:varchar(100);not null"`
	ISO3Code          *string    `json:"iso3_code,omitempty" gorm:"type:char(3);uniqueIndex"`
	NumericCode       *int16     `json:"numeric_code,omitempty" gorm:"type:smallint;uniqueIndex"`
	OfficialName      *string    `json:"official_name,omitempty" gorm:"type:varchar(200)"`
	CapitalCity       *string    `json:"capital_city,omitempty" gorm:"type:varchar(100)"`
	ContinentCode     *string    `json:"continent_code,omitempty" gorm:"type:varchar(2)"`
	RegionID          *uuid.UUID `json:"region_id,omitempty" gorm:"type:uuid"`
	PrimaryLanguageID *uuid.UUID `json:"primary_language_id,omitempty" gorm:"type:uuid"`
	CurrencyID        *uuid.UUID `json:"currency_id,omitempty" gorm:"type:uuid"`
	PhonePrefix       *string    `json:"phone_prefix,omitempty" gorm:"type:varchar(10)"`
	
	// Status Fields
	IsActive          bool       `json:"is_active" gorm:"default:true;not null"`
	IsDeleted         bool       `json:"is_deleted" gorm:"default:false;not null"`
	
	// LASANI Audit Fields (27 fields) - COMPLETE IMPLEMENTATION
	TenantID          string     `json:"tenant_id" gorm:"type:varchar(100);default:'default-tenant';index;not null"`
	CreatedAt         *time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;default:now()"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty" gorm:"type:uuid"`
	CreatedIP         *net.IP    `json:"created_ip,omitempty" gorm:"type:inet"`
	CreatedDevice     *json.RawMessage `json:"created_device,omitempty" gorm:"type:jsonb"`
	CreatedSession    *uuid.UUID `json:"created_session,omitempty" gorm:"type:uuid"`
	CreatedLocation   *json.RawMessage `json:"created_location,omitempty" gorm:"type:jsonb"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;default:now()"`
	UpdatedBy         *uuid.UUID `json:"updated_by,omitempty" gorm:"type:uuid"`
	UpdatedIP         *net.IP    `json:"updated_ip,omitempty" gorm:"type:inet"`
	UpdatedDevice     *json.RawMessage `json:"updated_device,omitempty" gorm:"type:jsonb"`
	UpdatedSession    *uuid.UUID `json:"updated_session,omitempty" gorm:"type:uuid"`
	UpdatedLocation   *json.RawMessage `json:"updated_location,omitempty" gorm:"type:jsonb"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty" gorm:"type:timestamptz"`
	DeletedBy         *uuid.UUID `json:"deleted_by,omitempty" gorm:"type:uuid"`
	DeletedIP         *net.IP    `json:"deleted_ip,omitempty" gorm:"type:inet"`
	DeletedDevice     *json.RawMessage `json:"deleted_device,omitempty" gorm:"type:jsonb"`
	DeletedSession    *uuid.UUID `json:"deleted_session,omitempty" gorm:"type:uuid"`
	DeletedLocation   *json.RawMessage `json:"deleted_location,omitempty" gorm:"type:jsonb"`
	SourceSystem      string     `json:"source_system" gorm:"type:varchar(50);default:'reference_master_geopolitical'"`
	ChangeReason      *string    `json:"change_reason,omitempty" gorm:"type:text"`
	Version           int        `json:"version" gorm:"default:1;not null"`
	DataClassification string    `json:"data_classification" gorm:"type:varchar(20);default:'PUBLIC'"`
	RetentionPolicy   *string    `json:"retention_policy,omitempty" gorm:"type:varchar(50)"`
	ComplianceFlags   *json.RawMessage `json:"compliance_flags,omitempty" gorm:"type:jsonb"`
	AuditTrail        *json.RawMessage `json:"audit_trail,omitempty" gorm:"type:jsonb"`
}

func (Country) TableName() string {
	return "domain_reference_master_geopolitical.countries"
}

// Region represents geographical regions
type Region struct {
	RegionID          uuid.UUID `json:"region_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RegionCode        string    `json:"region_code" gorm:"type:varchar(10);uniqueIndex;not null"`
	RegionName        string    `json:"region_name" gorm:"type:varchar(100);not null"`
	RegionType        string    `json:"region_type" gorm:"type:varchar(20);not null"`
	ParentRegionID    *uuid.UUID `json:"parent_region_id,omitempty" gorm:"type:uuid"`
	IsActive          bool      `json:"is_active" gorm:"default:true;not null"`
	IsDeleted         bool      `json:"is_deleted" gorm:"default:false;not null"`
	TenantID          string    `json:"tenant_id" gorm:"type:varchar(100);default:'default-tenant';index;not null"`
	CreatedAt         *time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;default:now()"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;default:now()"`
	Version           int       `json:"version" gorm:"default:1;not null"`
}

func (Region) TableName() string {
	return "domain_reference_master_geopolitical.regions"
}

// Language represents supported languages
type Language struct {
	LanguageID        uuid.UUID `json:"language_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LanguageCode      string    `json:"language_code" gorm:"type:char(2);uniqueIndex;not null"`
	LanguageName      string    `json:"language_name" gorm:"type:varchar(100);not null"`
	ISO3Code          *string   `json:"iso3_code,omitempty" gorm:"type:char(3);uniqueIndex"`
	NativeName        *string   `json:"native_name,omitempty" gorm:"type:varchar(100)"`
	Direction         string    `json:"direction" gorm:"type:varchar(3);default:'LTR';not null"`
	IsActive          bool      `json:"is_active" gorm:"default:true;not null"`
	IsDeleted         bool      `json:"is_deleted" gorm:"default:false;not null"`
	TenantID          string    `json:"tenant_id" gorm:"type:varchar(100);default:'default-tenant';index;not null"`
	CreatedAt         *time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;default:now()"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;default:now()"`
	Version           int       `json:"version" gorm:"default:1;not null"`
}

func (Language) TableName() string {
	return "domain_reference_master_geopolitical.languages"
}

// Timezone represents timezone information
type Timezone struct {
	TimezoneID        uuid.UUID `json:"timezone_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TimezoneCode      string    `json:"timezone_code" gorm:"type:varchar(10);uniqueIndex;not null"`
	TimezoneName      string    `json:"timezone_name" gorm:"type:varchar(100);not null"`
	UTCOffsetHours    int       `json:"utc_offset_hours" gorm:"not null"`
	UTCOffsetMinutes  int       `json:"utc_offset_minutes" gorm:"default:0;not null"`
	SupportsDST       bool      `json:"supports_dst" gorm:"default:false;not null"`
	DSTOffsetHours    *int      `json:"dst_offset_hours,omitempty"`
	IsActive          bool      `json:"is_active" gorm:"default:true;not null"`
	IsDeleted         bool      `json:"is_deleted" gorm:"default:false;not null"`
	TenantID          string    `json:"tenant_id" gorm:"type:varchar(100);default:'default-tenant';index;not null"`
	CreatedAt         *time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;default:now()"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;default:now()"`
	Version           int       `json:"version" gorm:"default:1;not null"`
}

func (Timezone) TableName() string {
	return "domain_reference_master_geopolitical.timezones"
}

// CountrySubdivision represents states, provinces, etc.
type CountrySubdivision struct {
	SubdivisionID     uuid.UUID `json:"subdivision_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SubdivisionCode   string    `json:"subdivision_code" gorm:"type:varchar(10);not null"`
	SubdivisionName   string    `json:"subdivision_name" gorm:"type:varchar(100);not null"`
	SubdivisionType   string    `json:"subdivision_type" gorm:"type:varchar(20);not null"`
	CountryID         uuid.UUID `json:"country_id" gorm:"type:uuid;not null"`
	IsActive          bool      `json:"is_active" gorm:"default:true;not null"`
	IsDeleted         bool      `json:"is_deleted" gorm:"default:false;not null"`
	TenantID          string    `json:"tenant_id" gorm:"type:varchar(100);default:'default-tenant';index;not null"`
	CreatedAt         *time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;default:now()"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;default:now()"`
	Version           int       `json:"version" gorm:"default:1;not null"`
}

func (CountrySubdivision) TableName() string {
	return "domain_reference_master_geopolitical.country_subdivisions"
}

// Locales represents locale information
type Locales struct {
	LocaleID          uuid.UUID `json:"locale_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LocaleCode        string    `json:"locale_code" gorm:"type:varchar(10);uniqueIndex;not null"`
	LocaleName        string    `json:"locale_name" gorm:"type:varchar(100);not null"`
	LanguageID        uuid.UUID `json:"language_id" gorm:"type:uuid;not null"`
	CountryID         *uuid.UUID `json:"country_id,omitempty" gorm:"type:uuid"`
	IsActive          bool      `json:"is_active" gorm:"default:true;not null"`
	IsDeleted         bool      `json:"is_deleted" gorm:"default:false;not null"`
	TenantID          string    `json:"tenant_id" gorm:"type:varchar(100);default:'default-tenant';index;not null"`
	CreatedAt         *time.Time `json:"created_at,omitempty" gorm:"type:timestamptz;default:now()"`
	UpdatedAt         *time.Time `json:"updated_at,omitempty" gorm:"type:timestamptz;default:now()"`
	Version           int       `json:"version" gorm:"default:1;not null"`
}

func (Locales) TableName() string {
	return "domain_reference_master_geopolitical.locales"
}