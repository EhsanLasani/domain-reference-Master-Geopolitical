-- ============================================================================
-- INDEXES: countries
-- PURPOSE: Performance indexes for countries entity only
-- ============================================================================

CREATE INDEX idx_countries_code ON domain_reference_master_geopolitical.countries(country_code);
CREATE INDEX idx_countries_active ON domain_reference_master_geopolitical.countries(is_active, is_deleted);
CREATE INDEX idx_countries_region_fk ON domain_reference_master_geopolitical.countries(region_id);
CREATE INDEX idx_countries_language_fk ON domain_reference_master_geopolitical.countries(primary_language_id);
CREATE INDEX idx_countries_continent ON domain_reference_master_geopolitical.countries(continent_code);