-- ============================================================================
-- INDEXES: regions
-- PURPOSE: Performance indexes for regions entity only
-- ============================================================================

CREATE INDEX idx_regions_code ON domain_reference_master_geopolitical.regions(region_code);
CREATE INDEX idx_regions_active ON domain_reference_master_geopolitical.regions(is_active, is_deleted);
CREATE INDEX idx_regions_parent ON domain_reference_master_geopolitical.regions(parent_region_id);