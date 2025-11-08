package models

import (
	"time"
	"github.com/google/uuid"
)

type Country struct {
	CountryID         uuid.UUID `json:"country_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CountryCode       string    `json:"country_code" gorm:"type:char(2);uniqueIndex;not null"`
	CountryName       string    `json:"country_name" gorm:"type:varchar(100);not null"`
	ISO3Code          string    `json:"iso3_code" gorm:"type:char(3);uniqueIndex"`
	NumericCode       int16     `json:"numeric_code" gorm:"type:smallint;uniqueIndex"`
	OfficialName      string    `json:"official_name" gorm:"type:varchar(200)"`
	CapitalCity       string    `json:"capital_city" gorm:"type:varchar(100)"`
	ContinentCode     string    `json:"continent_code" gorm:"type:varchar(2)"`
	RegionID          *uuid.UUID `json:"region_id,omitempty" gorm:"type:uuid"`
	PrimaryLanguageID *uuid.UUID `json:"primary_language_id,omitempty" gorm:"type:uuid"`
	CurrencyID        *uuid.UUID `json:"currency_id,omitempty" gorm:"type:uuid"`
	PhonePrefix       string    `json:"phone_prefix" gorm:"type:varchar(10)"`
	IsActive          bool      `json:"is_active" gorm:"default:true;not null"`
	IsDeleted         bool      `json:"is_deleted" gorm:"default:false;not null"`
	TenantID          string    `json:"tenant_id" gorm:"type:varchar(100);default:'default-tenant';index"`
	CreatedAt         time.Time `json:"created_at" gorm:"type:timestamptz;default:now()"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty" gorm:"type:uuid"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"type:timestamptz;default:now()"`
	UpdatedBy         *uuid.UUID `json:"updated_by,omitempty" gorm:"type:uuid"`
	Version           int       `json:"version" gorm:"default:1;not null"`
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