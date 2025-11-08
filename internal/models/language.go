package models

import (
	"time"
	"github.com/google/uuid"
)

type Language struct {
	LanguageID   uuid.UUID `json:"language_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	LanguageCode string    `json:"language_code" gorm:"type:char(2);uniqueIndex;not null"`
	LanguageName string    `json:"language_name" gorm:"type:varchar(100);not null"`
	ISO3Code     string    `json:"iso3_code" gorm:"type:char(3)"`
	NativeName   string    `json:"native_name" gorm:"type:varchar(100)"`
	Direction    string    `json:"direction" gorm:"type:varchar(3);default:'LTR'"`
	IsActive     bool      `json:"is_active" gorm:"default:true;not null"`
	IsDeleted    bool      `json:"is_deleted" gorm:"default:false;not null"`
	TenantID     string    `json:"tenant_id" gorm:"type:varchar(100);not null;index"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamptz;default:now()"`
	CreatedBy    *uuid.UUID `json:"created_by,omitempty" gorm:"type:uuid"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamptz;default:now()"`
	UpdatedBy    *uuid.UUID `json:"updated_by,omitempty" gorm:"type:uuid"`
	Version      int       `json:"version" gorm:"default:1;not null"`
}

func (Language) TableName() string {
	return "domain_reference_master_geopolitical.languages"
}

type LanguageInput struct {
	LanguageCode string `json:"language_code" binding:"required,len=2"`
	LanguageName string `json:"language_name" binding:"required,min=1,max=100"`
	ISO3Code     string `json:"iso3_code,omitempty" binding:"len=3"`
	NativeName   string `json:"native_name,omitempty" binding:"max=100"`
	Direction    string `json:"direction,omitempty" binding:"oneof=LTR RTL"`
}