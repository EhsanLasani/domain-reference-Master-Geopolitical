-- ============================================================================
-- ENTITY: Schema Setup
-- PURPOSE: Schema creation and enum definitions only
-- DEPENDENCIES: None
-- ============================================================================

CREATE SCHEMA IF NOT EXISTS domain_reference_master_geopolitical;
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
SET search_path TO domain_reference_master_geopolitical, lasani, public;

-- Continent enumeration for type safety
CREATE TYPE domain_reference_master_geopolitical.continent_enum AS ENUM (
    'AF', 'AS', 'EU', 'NA', 'SA', 'OC', 'AN'
);