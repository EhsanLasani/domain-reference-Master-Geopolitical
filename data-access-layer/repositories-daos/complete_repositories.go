package repositories

import (
	"context"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/orm-odm-abstractions"
)

// RegionRepository handles region CRUD operations
type RegionRepository struct {
	db *gorm.DB
}

func NewRegionRepository(db *gorm.DB) *RegionRepository {
	return &RegionRepository{db: db}
}

func (r *RegionRepository) GetAll(ctx context.Context, tenantID string) ([]models.Region, error) {
	var regions []models.Region
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND is_deleted = ?", tenantID, false).Find(&regions).Error
	return regions, err
}

func (r *RegionRepository) GetByCode(ctx context.Context, tenantID, code string) (*models.Region, error) {
	var region models.Region
	err := r.db.WithContext(ctx).Where("region_code = ? AND tenant_id = ? AND is_deleted = ?", code, tenantID, false).First(&region).Error
	if err == gorm.ErrRecordNotFound { return nil, nil }
	return &region, err
}

func (r *RegionRepository) Create(ctx context.Context, tenantID string, region *models.Region) error {
	now := time.Now()
	region.RegionID = uuid.New()
	region.TenantID = tenantID
	region.CreatedAt = &now
	region.UpdatedAt = &now
	region.IsActive = true
	region.IsDeleted = false
	region.Version = 1
	return r.db.WithContext(ctx).Create(region).Error
}

func (r *RegionRepository) Update(ctx context.Context, tenantID string, region *models.Region) error {
	now := time.Now()
	region.UpdatedAt = &now
	return r.db.WithContext(ctx).Where("region_code = ? AND tenant_id = ?", region.RegionCode, tenantID).Updates(region).Error
}

func (r *RegionRepository) Delete(ctx context.Context, tenantID, code string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Where("region_code = ? AND tenant_id = ?", code, tenantID).Updates(map[string]interface{}{
		"is_deleted": true, "deleted_at": now, "updated_at": now,
	}).Error
}

// LanguageRepository handles language CRUD operations
type LanguageRepository struct {
	db *gorm.DB
}

func NewLanguageRepository(db *gorm.DB) *LanguageRepository {
	return &LanguageRepository{db: db}
}

func (r *LanguageRepository) GetAll(ctx context.Context, tenantID string) ([]models.Language, error) {
	var languages []models.Language
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND is_deleted = ?", tenantID, false).Find(&languages).Error
	return languages, err
}

func (r *LanguageRepository) GetByCode(ctx context.Context, tenantID, code string) (*models.Language, error) {
	var language models.Language
	err := r.db.WithContext(ctx).Where("language_code = ? AND tenant_id = ? AND is_deleted = ?", code, tenantID, false).First(&language).Error
	if err == gorm.ErrRecordNotFound { return nil, nil }
	return &language, err
}

func (r *LanguageRepository) Create(ctx context.Context, tenantID string, language *models.Language) error {
	now := time.Now()
	language.LanguageID = uuid.New()
	language.TenantID = tenantID
	language.CreatedAt = &now
	language.UpdatedAt = &now
	language.IsActive = true
	language.IsDeleted = false
	language.Version = 1
	return r.db.WithContext(ctx).Create(language).Error
}

func (r *LanguageRepository) Update(ctx context.Context, tenantID string, language *models.Language) error {
	now := time.Now()
	language.UpdatedAt = &now
	return r.db.WithContext(ctx).Where("language_code = ? AND tenant_id = ?", language.LanguageCode, tenantID).Updates(language).Error
}

func (r *LanguageRepository) Delete(ctx context.Context, tenantID, code string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Where("language_code = ? AND tenant_id = ?", code, tenantID).Updates(map[string]interface{}{
		"is_deleted": true, "deleted_at": now, "updated_at": now,
	}).Error
}

// TimezoneRepository handles timezone CRUD operations
type TimezoneRepository struct {
	db *gorm.DB
}

func NewTimezoneRepository(db *gorm.DB) *TimezoneRepository {
	return &TimezoneRepository{db: db}
}

func (r *TimezoneRepository) GetAll(ctx context.Context, tenantID string) ([]models.Timezone, error) {
	var timezones []models.Timezone
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND is_deleted = ?", tenantID, false).Find(&timezones).Error
	return timezones, err
}

func (r *TimezoneRepository) GetByCode(ctx context.Context, tenantID, code string) (*models.Timezone, error) {
	var timezone models.Timezone
	err := r.db.WithContext(ctx).Where("timezone_code = ? AND tenant_id = ? AND is_deleted = ?", code, tenantID, false).First(&timezone).Error
	if err == gorm.ErrRecordNotFound { return nil, nil }
	return &timezone, err
}

func (r *TimezoneRepository) Create(ctx context.Context, tenantID string, timezone *models.Timezone) error {
	now := time.Now()
	timezone.TimezoneID = uuid.New()
	timezone.TenantID = tenantID
	timezone.CreatedAt = &now
	timezone.UpdatedAt = &now
	timezone.IsActive = true
	timezone.IsDeleted = false
	timezone.Version = 1
	return r.db.WithContext(ctx).Create(timezone).Error
}

func (r *TimezoneRepository) Update(ctx context.Context, tenantID string, timezone *models.Timezone) error {
	now := time.Now()
	timezone.UpdatedAt = &now
	return r.db.WithContext(ctx).Where("timezone_code = ? AND tenant_id = ?", timezone.TimezoneCode, tenantID).Updates(timezone).Error
}

func (r *TimezoneRepository) Delete(ctx context.Context, tenantID, code string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Where("timezone_code = ? AND tenant_id = ?", code, tenantID).Updates(map[string]interface{}{
		"is_deleted": true, "deleted_at": now, "updated_at": now,
	}).Error
}

// SubdivisionRepository handles subdivision CRUD operations
type SubdivisionRepository struct {
	db *gorm.DB
}

func NewSubdivisionRepository(db *gorm.DB) *SubdivisionRepository {
	return &SubdivisionRepository{db: db}
}

func (r *SubdivisionRepository) GetAll(ctx context.Context, tenantID string) ([]models.CountrySubdivision, error) {
	var subdivisions []models.CountrySubdivision
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND is_deleted = ?", tenantID, false).Find(&subdivisions).Error
	return subdivisions, err
}

func (r *SubdivisionRepository) GetByCountry(ctx context.Context, tenantID string, countryID uuid.UUID) ([]models.CountrySubdivision, error) {
	var subdivisions []models.CountrySubdivision
	err := r.db.WithContext(ctx).Where("country_id = ? AND tenant_id = ? AND is_deleted = ?", countryID, tenantID, false).Find(&subdivisions).Error
	return subdivisions, err
}

func (r *SubdivisionRepository) Create(ctx context.Context, tenantID string, subdivision *models.CountrySubdivision) error {
	now := time.Now()
	subdivision.SubdivisionID = uuid.New()
	subdivision.TenantID = tenantID
	subdivision.CreatedAt = &now
	subdivision.UpdatedAt = &now
	subdivision.IsActive = true
	subdivision.IsDeleted = false
	subdivision.Version = 1
	return r.db.WithContext(ctx).Create(subdivision).Error
}

func (r *SubdivisionRepository) Update(ctx context.Context, tenantID string, subdivision *models.CountrySubdivision) error {
	now := time.Now()
	subdivision.UpdatedAt = &now
	return r.db.WithContext(ctx).Where("subdivision_id = ? AND tenant_id = ?", subdivision.SubdivisionID, tenantID).Updates(subdivision).Error
}

func (r *SubdivisionRepository) Delete(ctx context.Context, tenantID string, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Where("subdivision_id = ? AND tenant_id = ?", id, tenantID).Updates(map[string]interface{}{
		"is_deleted": true, "deleted_at": now, "updated_at": now,
	}).Error
}

// LocaleRepository handles locale CRUD operations
type LocaleRepository struct {
	db *gorm.DB
}

func NewLocaleRepository(db *gorm.DB) *LocaleRepository {
	return &LocaleRepository{db: db}
}

func (r *LocaleRepository) GetAll(ctx context.Context, tenantID string) ([]models.Locales, error) {
	var locales []models.Locales
	err := r.db.WithContext(ctx).Where("tenant_id = ? AND is_deleted = ?", tenantID, false).Find(&locales).Error
	return locales, err
}

func (r *LocaleRepository) GetByCode(ctx context.Context, tenantID, code string) (*models.Locales, error) {
	var locale models.Locales
	err := r.db.WithContext(ctx).Where("locale_code = ? AND tenant_id = ? AND is_deleted = ?", code, tenantID, false).First(&locale).Error
	if err == gorm.ErrRecordNotFound { return nil, nil }
	return &locale, err
}

func (r *LocaleRepository) Create(ctx context.Context, tenantID string, locale *models.Locales) error {
	now := time.Now()
	locale.LocaleID = uuid.New()
	locale.TenantID = tenantID
	locale.CreatedAt = &now
	locale.UpdatedAt = &now
	locale.IsActive = true
	locale.IsDeleted = false
	locale.Version = 1
	return r.db.WithContext(ctx).Create(locale).Error
}

func (r *LocaleRepository) Update(ctx context.Context, tenantID string, locale *models.Locales) error {
	now := time.Now()
	locale.UpdatedAt = &now
	return r.db.WithContext(ctx).Where("locale_code = ? AND tenant_id = ?", locale.LocaleCode, tenantID).Updates(locale).Error
}

func (r *LocaleRepository) Delete(ctx context.Context, tenantID, code string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Where("locale_code = ? AND tenant_id = ?", code, tenantID).Updates(map[string]interface{}{
		"is_deleted": true, "deleted_at": now, "updated_at": now,
	}).Error
}