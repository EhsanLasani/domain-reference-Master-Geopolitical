-- ============================================================================
-- SCHEMA INITIALIZATION SCRIPT
-- PURPOSE: Create schema and initialize database for aligned application
-- DATABASE: geopolitical
-- SCHEMA: domain_reference_master_geopolitical
-- ============================================================================

-- Create schema
CREATE SCHEMA IF NOT EXISTS domain_reference_master_geopolitical;

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Set search path
SET search_path TO domain_reference_master_geopolitical, public;

-- Create continent enumeration
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'continent_enum') THEN
        CREATE TYPE domain_reference_master_geopolitical.continent_enum AS ENUM (
            'AF', 'AS', 'EU', 'NA', 'SA', 'OC', 'AN'
        );
    END IF;
END $$;

-- Grant permissions
GRANT USAGE ON SCHEMA domain_reference_master_geopolitical TO postgres;
GRANT ALL PRIVILEGES ON SCHEMA domain_reference_master_geopolitical TO postgres;

-- Confirm schema creation
SELECT 'Schema domain_reference_master_geopolitical created successfully' as status;