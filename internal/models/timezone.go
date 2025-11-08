package models

import (
	"time"
	"github.com/google/uuid"
)

type Timezone struct {
	TimezoneID       uuid.UUID `json:"timezone_id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	TimezoneCode     string    `json:"timezone_code" gorm:"type:varchar(50);uniqueIndex;not null"`
	TimezoneName     string    `json:"timezone_name" gorm:"type:varchar(100);not null"`
	UTCOffsetHours   int16     `json:"utc_offset_hours" gorm:"not null"`
	UTCOffsetMinutes int16     `json:"utc_offset_minutes" gorm:"default:0"`
	SupportsDST      bool      `json:"supports_dst" gorm:"default:false"`
	DSTOffsetHours   *int16    `json:"dst_offset_hours,omitempty" gorm:"type:smallint"`
	IsActive         bool      `json:"is_active" gorm:"default:true;not null"`
	IsDeleted        bool      `json:"is_deleted" gorm:"default:false;not null"`
	TenantID         string    `json:"tenant_id" gorm:"type:varchar(100);not null;index"`
	CreatedAt        time.Time `json:"created_at" gorm:"type:timestamptz;default:now()"`
	CreatedBy        *uuid.UUID `json:"created_by,omitempty" gorm:"type:uuid"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"type:timestamptz;default:now()"`
	UpdatedBy        *uuid.UUID `json:"updated_by,omitempty" gorm:"type:uuid"`
	Version          int       `json:"version" gorm:"default:1;not null"`
}

func (Timezone) TableName() string {
	return "domain_reference_master_geopolitical.timezones"
}

type TimezoneInput struct {
	TimezoneCode     string `json:"timezone_code" binding:"required,min=1,max=50"`
	TimezoneName     string `json:"timezone_name" binding:"required,min=1,max=100"`
	UTCOffsetHours   int16  `json:"utc_offset_hours" binding:"required,min=-12,max=14"`
	UTCOffsetMinutes int16  `json:"utc_offset_minutes,omitempty" binding:"min=0,max=59"`
	SupportsDST      bool   `json:"supports_dst,omitempty"`
	DSTOffsetHours   *int16 `json:"dst_offset_hours,omitempty" binding:"omitempty,min=-12,max=14"`
}