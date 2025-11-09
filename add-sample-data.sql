-- Add sample countries
INSERT INTO domain_reference_master_geopolitical.countries 
(country_id, country_code, country_name, iso3_code, numeric_code, official_name, capital_city, continent_code, tenant_id, is_active, is_deleted, created_at, updated_at, version, source_system)
VALUES 
(gen_random_uuid(), 'US', 'United States', 'USA', 840, 'United States of America', 'Washington D.C.', 'NA', 'default-tenant', true, false, NOW(), NOW(), 1, 'reference_master_geopolitical'),
(gen_random_uuid(), 'GB', 'United Kingdom', 'GBR', 826, 'United Kingdom of Great Britain and Northern Ireland', 'London', 'EU', 'default-tenant', true, false, NOW(), NOW(), 1, 'reference_master_geopolitical'),
(gen_random_uuid(), 'FR', 'France', 'FRA', 250, 'French Republic', 'Paris', 'EU', 'default-tenant', true, false, NOW(), NOW(), 1, 'reference_master_geopolitical'),
(gen_random_uuid(), 'DE', 'Germany', 'DEU', 276, 'Federal Republic of Germany', 'Berlin', 'EU', 'default-tenant', true, false, NOW(), NOW(), 1, 'reference_master_geopolitical'),
(gen_random_uuid(), 'JP', 'Japan', 'JPN', 392, 'Japan', 'Tokyo', 'AS', 'default-tenant', true, false, NOW(), NOW(), 1, 'reference_master_geopolitical')
ON CONFLICT (country_code) DO NOTHING;