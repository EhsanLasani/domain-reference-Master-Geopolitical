package models

import (
	"time"
	"github.com/google/uuid"
)

type CountrySubdivision struct {
	SubdivisionID       uuid.UUID `json:"subdivision_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	SubdivisionCode     string    `json:"subdivision_code" gorm:"type:varchar(10);not null"`
	SubdivisionName     string    `json:"subdivision_name" gorm:"type:varchar(100);not null"`
	CountryID           uuid.UUID `json:"country_id" gorm:"type:uuid;not null"`
	SubdivisionType     string    `json:"subdivision_type" gorm:"type:varchar(20)"`
	ParentSubdivisionID *uuid.UUID `json:"parent_subdivision_id,omitempty" gorm:"type:uuid"`
	IsActive            bool      `json:"is_active" gorm:"default:true;not null"`
	IsDeleted           bool      `json:"is_deleted" gorm:"default:false;not null"`
	TenantID            string    `json:"tenant_id" gorm:"type:varchar(100);not null;index"`
	CreatedAt           time.Time `json:"created_at" gorm:"type:timestamptz;default:now()"`
	CreatedBy           *uuid.UUID `json:"created_by,omitempty" gorm:"type:uuid"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"type:timestamptz;default:now()"`
	UpdatedBy           *uuid.UUID `json:"updated_by,omitempty" gorm:"type:uuid"`
	Version             int       `json:"version" gorm:"default:1;not null"`
}

func (CountrySubdivision) TableName() string {
	return "domain_reference_master_geopolitical.country_subdivisions"
}

type SubdivisionInput struct {
	SubdivisionCode     string     `json:"subdivision_code" binding:"required,min=1,max=10"`
	SubdivisionName     string     `json:"subdivision_name" binding:"required,min=1,max=100"`
	CountryID           uuid.UUID  `json:"country_id" binding:"required"`
	SubdivisionType     string     `json:"subdivision_type,omitempty" binding:"max=20"`
	ParentSubdivisionID *uuid.UUID `json:"parent_subdivision_id,omitempty"`
}