package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/models"
)

type RegionRepository struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) *RegionRepository {
	return &RegionRepository{db: db}
}

func (r *RegionRepository) Create(ctx context.Context, tenantID string, region *models.Region) error {
	region.TenantID = tenantID
	return r.db.WithContext(ctx).Create(region).Error
}

func (r *RegionRepository) GetByID(ctx context.Context, tenantID string, id uuid.UUID) (*models.Region, error) {
	var region models.Region
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND region_id = ? AND is_deleted = false", tenantID, id).
		First(&region).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *RegionRepository) GetByCode(ctx context.Context, tenantID string, code string) (*models.Region, error) {
	var region models.Region
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND region_code = ? AND is_deleted = false", tenantID, code).
		First(&region).Error
	if err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *RegionRepository) List(ctx context.Context, tenantID string, limit, offset int) ([]*models.Region, int64, error) {
	var regions []*models.Region
	var total int64
	
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_deleted = false", tenantID)
	
	if err := query.Model(&models.Region{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := query.Limit(limit).Offset(offset).Find(&regions).Error
	return regions, total, err
}

func (r *RegionRepository) Update(ctx context.Context, tenantID string, region *models.Region) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND region_id = ?", tenantID, region.RegionID).
		Updates(region).Error
}

func (r *RegionRepository) Delete(ctx context.Context, tenantID string, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Region{}).
		Where("tenant_id = ? AND region_id = ?", tenantID, id).
		Update("is_deleted", true).Error
}