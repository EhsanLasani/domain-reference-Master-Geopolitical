-- ============================================================================
-- FILE: database_views.sql
-- DOMAIN: Reference Master Geopolitical
-- LAYER: Data Access Layer - Views
-- PURPOSE: Database views for optimized queries and business operations
-- VERSION: 1.0.0
-- CREATED: 2025-11-07
-- ============================================================================

-- Active Countries View - Most frequently used
CREATE OR REPLACE VIEW domain_reference_master_geopolitical.v_countries_active AS
SELECT 
    country_id,
    country_code,
    country_name,
    iso3_code,
    official_name,
    capital_city,
    continent_code,
    phone_prefix,
    created_at,
    updated_at,
    version
FROM domain_reference_master_geopolitical.countries
WHERE is_active = true AND is_deleted = false
ORDER BY country_code;

-- Countries with Regions View - For reporting
CREATE OR REPLACE VIEW domain_reference_master_geopolitical.v_countries_with_regions AS
SELECT 
    c.country_id,
    c.country_code,
    c.country_name,
    c.iso3_code,
    c.continent_code,
    r.region_id,
    r.region_code,
    r.region_name,
    r.region_type
FROM domain_reference_master_geopolitical.countries c
LEFT JOIN domain_reference_master_geopolitical.regions r ON c.region_id = r.region_id
WHERE c.is_active = true AND c.is_deleted = false
ORDER BY c.country_code;

-- Audit Trail View - For compliance
CREATE OR REPLACE VIEW domain_reference_master_geopolitical.v_audit_trail AS
SELECT 
    'countries' as entity_type,
    country_id as entity_id,
    country_code as entity_code,
    created_at,
    created_by,
    updated_at,
    updated_by,
    version,
    change_reason
FROM domain_reference_master_geopolitical.countries
ORDER BY updated_at DESC;