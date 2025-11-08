package applicationservices

import (
	"context"
	"github.com/google/uuid"
	repositories "github.com/EhsanLasani/domain-reference-Master-Geopolitical/data-access-layer/repositories-daos"
	models "github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/models"
	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/logging"
)

type RegionAppService struct {
	regionRepo *repositories.RegionRepository
	logger     *logging.StructuredLogger
}

func NewRegionAppService(regionRepo *repositories.RegionRepository, logger *logging.StructuredLogger) *RegionAppService {
	return &RegionAppService{
		regionRepo: regionRepo,
		logger:     logger,
	}
}

func (s *RegionAppService) CreateRegion(ctx context.Context, tenantID string, input *models.RegionInput) (*models.Region, error) {
	s.logger.Info(ctx, "Creating region", logging.Field{Key: "region_code", Value: input.RegionCode})
	
	region := &models.Region{
		RegionID:       uuid.New(),
		RegionCode:     input.RegionCode,
		RegionName:     input.RegionName,
		RegionType:     input.RegionType,
		ParentRegionID: input.ParentRegionID,
	}
	
	if err := s.regionRepo.Create(ctx, tenantID, region); err != nil {
		s.logger.Error(ctx, "Failed to create region", err)
		return nil, err
	}
	
	return region, nil
}

func (s *RegionAppService) GetRegion(ctx context.Context, tenantID string, id uuid.UUID) (*models.Region, error) {
	return s.regionRepo.GetByID(ctx, tenantID, id)
}

func (s *RegionAppService) GetRegionByCode(ctx context.Context, tenantID string, code string) (*models.Region, error) {
	return s.regionRepo.GetByCode(ctx, tenantID, code)
}

func (s *RegionAppService) ListRegions(ctx context.Context, tenantID string, limit, offset int) ([]*models.Region, int64, error) {
	return s.regionRepo.List(ctx, tenantID, limit, offset)
}

func (s *RegionAppService) UpdateRegion(ctx context.Context, tenantID string, id uuid.UUID, input *models.RegionInput) (*models.Region, error) {
	region, err := s.regionRepo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	
	region.RegionName = input.RegionName
	region.RegionType = input.RegionType
	region.ParentRegionID = input.ParentRegionID
	region.Version++
	
	if err := s.regionRepo.Update(ctx, tenantID, region); err != nil {
		return nil, err
	}
	
	return region, nil
}

func (s *RegionAppService) DeleteRegion(ctx context.Context, tenantID string, id uuid.UUID) error {
	return s.regionRepo.Delete(ctx, tenantID, id)
}