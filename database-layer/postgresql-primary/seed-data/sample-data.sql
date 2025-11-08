-- ============================================================================
-- SAMPLE DATA SEEDING
-- PURPOSE: Insert sample data for all geopolitical entities
-- SCHEMA: domain_reference_master_geopolitical
-- ============================================================================

SET search_path TO domain_reference_master_geopolitical, public;

-- Sample Countries
INSERT INTO countries (
    country_id, country_code, country_name, iso3_code, numeric_code, 
    official_name, capital_city, continent_code, phone_prefix, 
    is_active, is_deleted, tenant_id, created_at, created_by, updated_at, updated_by, version
) VALUES 
(gen_random_uuid(), 'US', 'United States', 'USA', 840, 'United States of America', 'Washington D.C.', 'NA', '+1', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'GB', 'United Kingdom', 'GBR', 826, 'United Kingdom of Great Britain and Northern Ireland', 'London', 'EU', '+44', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'DE', 'Germany', 'DEU', 276, 'Federal Republic of Germany', 'Berlin', 'EU', '+49', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'FR', 'France', 'FRA', 250, 'French Republic', 'Paris', 'EU', '+33', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'JP', 'Japan', 'JPN', 392, 'Japan', 'Tokyo', 'AS', '+81', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'CN', 'China', 'CHN', 156, 'People''s Republic of China', 'Beijing', 'AS', '+86', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'IN', 'India', 'IND', 356, 'Republic of India', 'New Delhi', 'AS', '+91', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'BR', 'Brazil', 'BRA', 76, 'Federative Republic of Brazil', 'Brasília', 'SA', '+55', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'AU', 'Australia', 'AUS', 36, 'Commonwealth of Australia', 'Canberra', 'OC', '+61', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'CA', 'Canada', 'CAN', 124, 'Canada', 'Ottawa', 'NA', '+1', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1);

-- Sample Regions
INSERT INTO regions (
    region_id, region_code, region_name, region_type, parent_region_id,
    is_active, is_deleted, tenant_id, created_at, created_by, updated_at, updated_by, version
) VALUES 
(gen_random_uuid(), 'NA', 'North America', 'CONTINENT', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'EU', 'Europe', 'CONTINENT', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'AS', 'Asia', 'CONTINENT', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'SA', 'South America', 'CONTINENT', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'AF', 'Africa', 'CONTINENT', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'OC', 'Oceania', 'CONTINENT', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'WE', 'Western Europe', 'REGION', (SELECT region_id FROM regions WHERE region_code = 'EU'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'EE', 'Eastern Europe', 'REGION', (SELECT region_id FROM regions WHERE region_code = 'EU'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'SEA', 'Southeast Asia', 'REGION', (SELECT region_id FROM regions WHERE region_code = 'AS'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'ME', 'Middle East', 'REGION', (SELECT region_id FROM regions WHERE region_code = 'AS'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1);

-- Sample Languages
INSERT INTO languages (
    language_id, language_code, language_name, iso3_code, native_name, direction,
    is_active, is_deleted, tenant_id, created_at, created_by, updated_at, updated_by, version
) VALUES 
(gen_random_uuid(), 'en', 'English', 'eng', 'English', 'LTR', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'es', 'Spanish', 'spa', 'Español', 'LTR', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'fr', 'French', 'fra', 'Français', 'LTR', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'de', 'German', 'deu', 'Deutsch', 'LTR', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'zh', 'Chinese', 'zho', '中文', 'LTR', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'ja', 'Japanese', 'jpn', '日本語', 'LTR', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'ar', 'Arabic', 'ara', 'العربية', 'RTL', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'hi', 'Hindi', 'hin', 'हिन्दी', 'LTR', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'pt', 'Portuguese', 'por', 'Português', 'LTR', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'ru', 'Russian', 'rus', 'Русский', 'LTR', true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1);

-- Sample Timezones
INSERT INTO timezones (
    timezone_id, timezone_code, timezone_name, utc_offset_hours, utc_offset_minutes, 
    supports_dst, dst_offset_hours, is_active, is_deleted, tenant_id, 
    created_at, created_by, updated_at, updated_by, version
) VALUES 
(gen_random_uuid(), 'UTC', 'Coordinated Universal Time', 0, 0, false, NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'EST', 'Eastern Standard Time', -5, 0, true, 1, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'PST', 'Pacific Standard Time', -8, 0, true, 1, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'GMT', 'Greenwich Mean Time', 0, 0, true, 1, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'CET', 'Central European Time', 1, 0, true, 1, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'JST', 'Japan Standard Time', 9, 0, false, NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'CST', 'China Standard Time', 8, 0, false, NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'IST', 'India Standard Time', 5, 30, false, NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'AEST', 'Australian Eastern Standard Time', 10, 0, true, 1, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'BRT', 'Brasília Time', -3, 0, false, NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1);

-- Sample Country Subdivisions
INSERT INTO country_subdivisions (
    subdivision_id, subdivision_code, subdivision_name, country_id, subdivision_type, parent_subdivision_id,
    is_active, is_deleted, tenant_id, created_at, created_by, updated_at, updated_by, version
) VALUES 
(gen_random_uuid(), 'US-CA', 'California', (SELECT country_id FROM countries WHERE country_code = 'US'), 'STATE', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'US-NY', 'New York', (SELECT country_id FROM countries WHERE country_code = 'US'), 'STATE', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'US-TX', 'Texas', (SELECT country_id FROM countries WHERE country_code = 'US'), 'STATE', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'GB-ENG', 'England', (SELECT country_id FROM countries WHERE country_code = 'GB'), 'COUNTRY', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'GB-SCT', 'Scotland', (SELECT country_id FROM countries WHERE country_code = 'GB'), 'COUNTRY', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'DE-BY', 'Bavaria', (SELECT country_id FROM countries WHERE country_code = 'DE'), 'STATE', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'FR-IDF', 'Île-de-France', (SELECT country_id FROM countries WHERE country_code = 'FR'), 'REGION', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'JP-13', 'Tokyo', (SELECT country_id FROM countries WHERE country_code = 'JP'), 'PREFECTURE', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'CN-BJ', 'Beijing', (SELECT country_id FROM countries WHERE country_code = 'CN'), 'MUNICIPALITY', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'IN-DL', 'Delhi', (SELECT country_id FROM countries WHERE country_code = 'IN'), 'TERRITORY', NULL, true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1);

-- Sample Locales
INSERT INTO locales (
    locale_id, locale_code, locale_name, language_id, country_id,
    is_active, is_deleted, tenant_id, created_at, created_by, updated_at, updated_by, version
) VALUES 
(gen_random_uuid(), 'en-US', 'English (United States)', (SELECT language_id FROM languages WHERE language_code = 'en'), (SELECT country_id FROM countries WHERE country_code = 'US'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'en-GB', 'English (United Kingdom)', (SELECT language_id FROM languages WHERE language_code = 'en'), (SELECT country_id FROM countries WHERE country_code = 'GB'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'fr-FR', 'French (France)', (SELECT language_id FROM languages WHERE language_code = 'fr'), (SELECT country_id FROM countries WHERE country_code = 'FR'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'de-DE', 'German (Germany)', (SELECT language_id FROM languages WHERE language_code = 'de'), (SELECT country_id FROM countries WHERE country_code = 'DE'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'ja-JP', 'Japanese (Japan)', (SELECT language_id FROM languages WHERE language_code = 'ja'), (SELECT country_id FROM countries WHERE country_code = 'JP'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'zh-CN', 'Chinese (China)', (SELECT language_id FROM languages WHERE language_code = 'zh'), (SELECT country_id FROM countries WHERE country_code = 'CN'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'hi-IN', 'Hindi (India)', (SELECT language_id FROM languages WHERE language_code = 'hi'), (SELECT country_id FROM countries WHERE country_code = 'IN'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'pt-BR', 'Portuguese (Brazil)', (SELECT language_id FROM languages WHERE language_code = 'pt'), (SELECT country_id FROM countries WHERE country_code = 'BR'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'en-AU', 'English (Australia)', (SELECT language_id FROM languages WHERE language_code = 'en'), (SELECT country_id FROM countries WHERE country_code = 'AU'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1),
(gen_random_uuid(), 'en-CA', 'English (Canada)', (SELECT language_id FROM languages WHERE language_code = 'en'), (SELECT country_id FROM countries WHERE country_code = 'CA'), true, false, 'default-tenant', NOW(), 'system', NOW(), 'system', 1);

-- Verify data insertion
SELECT 'Countries' as entity, COUNT(*) as count FROM countries
UNION ALL
SELECT 'Regions' as entity, COUNT(*) as count FROM regions
UNION ALL
SELECT 'Languages' as entity, COUNT(*) as count FROM languages
UNION ALL
SELECT 'Timezones' as entity, COUNT(*) as count FROM timezones
UNION ALL
SELECT 'Subdivisions' as entity, COUNT(*) as count FROM country_subdivisions
UNION ALL
SELECT 'Locales' as entity, COUNT(*) as count FROM locales;