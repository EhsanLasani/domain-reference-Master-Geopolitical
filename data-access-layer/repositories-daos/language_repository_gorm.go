package repositories

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/models"
)

type LanguageRepository struct {
	db *gorm.DB
}

func NewLanguageRepository(db *gorm.DB) *LanguageRepository {
	return &LanguageRepository{db: db}
}

func (r *LanguageRepository) Create(ctx context.Context, tenantID string, language *models.Language) error {
	language.TenantID = tenantID
	return r.db.WithContext(ctx).Create(language).Error
}

func (r *LanguageRepository) GetByID(ctx context.Context, tenantID string, id uuid.UUID) (*models.Language, error) {
	var language models.Language
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND language_id = ? AND is_deleted = false", tenantID, id).
		First(&language).Error
	if err != nil {
		return nil, err
	}
	return &language, nil
}

func (r *LanguageRepository) GetByCode(ctx context.Context, tenantID string, code string) (*models.Language, error) {
	var language models.Language
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND language_code = ? AND is_deleted = false", tenantID, code).
		First(&language).Error
	if err != nil {
		return nil, err
	}
	return &language, nil
}

func (r *LanguageRepository) List(ctx context.Context, tenantID string, limit, offset int) ([]*models.Language, int64, error) {
	var languages []*models.Language
	var total int64
	
	query := r.db.WithContext(ctx).
		Where("tenant_id = ? AND is_deleted = false", tenantID)
	
	if err := query.Model(&models.Language{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	err := query.Limit(limit).Offset(offset).Find(&languages).Error
	return languages, total, err
}

func (r *LanguageRepository) Update(ctx context.Context, tenantID string, language *models.Language) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ? AND language_id = ?", tenantID, language.LanguageID).
		Updates(language).Error
}

func (r *LanguageRepository) Delete(ctx context.Context, tenantID string, id uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&models.Language{}).
		Where("tenant_id = ? AND language_id = ?", tenantID, id).
		Update("is_deleted", true).Error
}