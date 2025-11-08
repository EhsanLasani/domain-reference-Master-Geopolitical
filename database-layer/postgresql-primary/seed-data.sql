-- ============================================================================
-- REFERENCE MASTER GEOPOLITICAL - SEED DATA
-- ============================================================================
-- Essential reference data for development and testing
-- ============================================================================

-- System user for seeding
DO $$
DECLARE
    system_user_id UUID := '00000000-0000-0000-0000-000000000001';
    seed_timestamp TIMESTAMPTZ := now();
BEGIN

-- Seed Regions (Continents)
INSERT INTO domain_reference_master_geopolitical.regions (
    region_code, region_name, region_type,
    is_active, is_deleted, created_at, created_by, source_system, change_reason, version
) VALUES 
    ('AF', 'Africa', 'Continent', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('AS', 'Asia', 'Continent', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('EU', 'Europe', 'Continent', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('NA', 'North America', 'Continent', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('SA', 'South America', 'Continent', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('OC', 'Oceania', 'Continent', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('AN', 'Antarctica', 'Continent', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1);

-- Seed Major Languages
INSERT INTO domain_reference_master_geopolitical.languages (
    language_code, language_name, iso3_code, native_name, direction,
    is_active, is_deleted, created_at, created_by, source_system, change_reason, version
) VALUES 
    ('en', 'English', 'eng', 'English', 'LTR', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('es', 'Spanish', 'spa', 'Español', 'LTR', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('fr', 'French', 'fra', 'Français', 'LTR', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('de', 'German', 'deu', 'Deutsch', 'LTR', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('zh', 'Chinese', 'zho', '中文', 'LTR', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('ja', 'Japanese', 'jpn', '日本語', 'LTR', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('ar', 'Arabic', 'ara', 'العربية', 'RTL', true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1);

-- Seed Major Timezones
INSERT INTO domain_reference_master_geopolitical.timezones (
    timezone_code, timezone_name, utc_offset_hours, utc_offset_minutes, supports_dst,
    is_active, is_deleted, created_at, created_by, source_system, change_reason, version
) VALUES 
    ('UTC', 'Coordinated Universal Time', 0, 0, false, true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('America/New_York', 'Eastern Time', -5, 0, true, true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('America/Los_Angeles', 'Pacific Time', -8, 0, true, true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('Europe/London', 'Greenwich Mean Time', 0, 0, true, true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('Europe/Paris', 'Central European Time', 1, 0, true, true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('Asia/Tokyo', 'Japan Standard Time', 9, 0, false, true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1);

-- Seed Major Countries
INSERT INTO domain_reference_master_geopolitical.countries (
    country_code, country_name, iso3_code, numeric_code, official_name, 
    capital_city, continent_code, currency_id, phone_prefix,
    region_id, primary_language_id,
    is_active, is_deleted, created_at, created_by, source_system, change_reason, version
) VALUES 
    ('US', 'United States', 'USA', 840, 'United States of America', 'Washington, D.C.', 'NA', null, '+1',
     (SELECT region_id FROM domain_reference_master_geopolitical.regions WHERE region_code = 'NA'),
     (SELECT language_id FROM domain_reference_master_geopolitical.languages WHERE language_code = 'en'),
     true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('GB', 'United Kingdom', 'GBR', 826, 'United Kingdom of Great Britain and Northern Ireland', 'London', 'EU', null, '+44',
     (SELECT region_id FROM domain_reference_master_geopolitical.regions WHERE region_code = 'EU'),
     (SELECT language_id FROM domain_reference_master_geopolitical.languages WHERE language_code = 'en'),
     true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('FR', 'France', 'FRA', 250, 'French Republic', 'Paris', 'EU', null, '+33',
     (SELECT region_id FROM domain_reference_master_geopolitical.regions WHERE region_code = 'EU'),
     (SELECT language_id FROM domain_reference_master_geopolitical.languages WHERE language_code = 'fr'),
     true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('DE', 'Germany', 'DEU', 276, 'Federal Republic of Germany', 'Berlin', 'EU', null, '+49',
     (SELECT region_id FROM domain_reference_master_geopolitical.regions WHERE region_code = 'EU'),
     (SELECT language_id FROM domain_reference_master_geopolitical.languages WHERE language_code = 'de'),
     true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('JP', 'Japan', 'JPN', 392, 'Japan', 'Tokyo', 'AS', null, '+81',
     (SELECT region_id FROM domain_reference_master_geopolitical.regions WHERE region_code = 'AS'),
     (SELECT language_id FROM domain_reference_master_geopolitical.languages WHERE language_code = 'ja'),
     true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1);

-- Seed Major Locales
INSERT INTO domain_reference_master_geopolitical.locales (
    locale_code, locale_name, language_id, country_id,
    is_active, is_deleted, created_at, created_by, source_system, change_reason, version
) VALUES 
    ('en-US', 'English (United States)', 
     (SELECT language_id FROM domain_reference_master_geopolitical.languages WHERE language_code = 'en'), 
     (SELECT country_id FROM domain_reference_master_geopolitical.countries WHERE country_code = 'US'), 
     true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('en-GB', 'English (United Kingdom)', 
     (SELECT language_id FROM domain_reference_master_geopolitical.languages WHERE language_code = 'en'), 
     (SELECT country_id FROM domain_reference_master_geopolitical.countries WHERE country_code = 'GB'), 
     true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1),
    ('fr-FR', 'French (France)', 
     (SELECT language_id FROM domain_reference_master_geopolitical.languages WHERE language_code = 'fr'), 
     (SELECT country_id FROM domain_reference_master_geopolitical.countries WHERE country_code = 'FR'), 
     true, false, seed_timestamp, system_user_id, 'seed_data', 'Initial seed data', 1);

END $$;