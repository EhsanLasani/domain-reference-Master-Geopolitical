package models

import (
	"time"
	"github.com/google/uuid"
)

type Region struct {
	RegionID uuid.UUID `json:"region_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	RegionCode string `json:"region_code" gorm:"type:varchar(10);uniqueIndex;not null"`
	RegionName string `json:"region_name" gorm:"type:varchar(100);not null"`
	RegionType string `json:"region_type" gorm:"type:varchar(20)"`
	ParentRegionID *uuid.UUID `json:"parent_region_id,omitempty" gorm:"type:uuid"`
	IsActive  bool `json:"is_active" gorm:"default:true;not null"`
	IsDeleted bool `json:"is_deleted" gorm:"default:false;not null"`
	TenantID        string     `json:"tenant_id" gorm:"type:varchar(100);not null;index"`
	CreatedAt       time.Time  `json:"created_at" gorm:"type:timestamptz;default:now()"`
	CreatedBy       *uuid.UUID `json:"created_by,omitempty" gorm:"type:uuid"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"type:timestamptz;default:now()"`
	UpdatedBy       *uuid.UUID `json:"updated_by,omitempty" gorm:"type:uuid"`
	Version         int        `json:"version" gorm:"default:1;not null"`
}

func (Region) TableName() string {
	return "domain_reference_master_geopolitical.regions"
}

type RegionInput struct {
	RegionCode     string     `json:"region_code" binding:"required,min=1,max=10"`
	RegionName     string     `json:"region_name" binding:"required,min=1,max=100"`
	RegionType     string     `json:"region_type,omitempty" binding:"max=20"`
	ParentRegionID *uuid.UUID `json:"parent_region_id,omitempty"`
}