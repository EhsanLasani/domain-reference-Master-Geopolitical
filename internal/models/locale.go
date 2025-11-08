package models

import (
	"time"
	"github.com/google/uuid"
)

type Locales struct {
	LocaleID   uuid.UUID `json:"locale_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LocaleCode string    `json:"locale_code" gorm:"type:varchar(10);uniqueIndex;not null"`
	LocaleName string    `json:"locale_name" gorm:"type:varchar(100);not null"`
	LanguageID uuid.UUID `json:"language_id" gorm:"type:uuid;not null"`
	// Language relationship removed to prevent circular FK constraint
	CountryID  *uuid.UUID `json:"country_id,omitempty" gorm:"type:uuid"`
	// Country relationship removed to prevent circular FK constraint
	IsActive   bool      `json:"is_active" gorm:"default:true;not null"`
	IsDeleted  bool      `json:"is_deleted" gorm:"default:false;not null"`
	TenantID   string    `json:"tenant_id" gorm:"type:varchar(100);not null;index"`
	CreatedAt  time.Time `json:"created_at" gorm:"type:timestamptz;default:now()"`
	CreatedBy  *uuid.UUID `json:"created_by,omitempty" gorm:"type:uuid"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:timestamptz;default:now()"`
	UpdatedBy  *uuid.UUID `json:"updated_by,omitempty" gorm:"type:uuid"`
	Version    int       `json:"version" gorm:"default:1;not null"`
}

func (Locales) TableName() string {
	return "domain_reference_master_geopolitical.locales"
}

type LocalesInput struct {
	LocaleCode string     `json:"locale_code" binding:"required,min=1,max=10"`
	LocaleName string     `json:"locale_name" binding:"required,min=1,max=100"`
	LanguageID uuid.UUID  `json:"language_id" binding:"required"`
	CountryID  *uuid.UUID `json:"country_id,omitempty"`
}