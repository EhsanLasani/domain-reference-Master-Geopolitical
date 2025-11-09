package models

import (
	"time"
	"net"
	"encoding/json"
	"github.com/google/uuid"
)

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
	
	// LASANI Audit Fields (27 fields)
	TenantID          string     `json:"tenant_id" gorm:"type:varchar(100);default:'default-tenant';index"`
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
}

func (Country) TableName() string {
	return "domain_reference_master_geopolitical.countries"
}

type CountryInput struct {
	CountryCode       string     `json:"country_code" binding:"required,len=2"`
	CountryName       string     `json:"country_name" binding:"required,min=1,max=100"`
	ISO3Code          string     `json:"iso3_code,omitempty" binding:"len=3"`
	NumericCode       int16      `json:"numeric_code,omitempty"`
	OfficialName      string     `json:"official_name,omitempty" binding:"max=200"`
	CapitalCity       string     `json:"capital_city,omitempty" binding:"max=100"`
	ContinentCode     string     `json:"continent_code,omitempty" binding:"len=2"`
	RegionID          *uuid.UUID `json:"region_id,omitempty"`
	PrimaryLanguageID *uuid.UUID `json:"primary_language_id,omitempty"`
	PhonePrefix       string     `json:"phone_prefix,omitempty" binding:"max=10"`
}