-- ============================================================================
-- ENTITY: countries
-- PURPOSE: Countries table with LASANI compliance
-- DEPENDENCIES: regions, languages
-- ============================================================================

CREATE TABLE domain_reference_master_geopolitical.countries (
    country_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code CHAR(2) COLLATE "C" UNIQUE NOT NULL,
    country_name VARCHAR(100) NOT NULL,
    iso3_code CHAR(3) COLLATE "C" UNIQUE,
    numeric_code SMALLINT UNIQUE,
    official_name VARCHAR(200),
    capital_city VARCHAR(100),
    continent_code domain_reference_master_geopolitical.continent_enum,
    region_id UUID REFERENCES domain_reference_master_geopolitical.regions(region_id),
    primary_language_id UUID REFERENCES domain_reference_master_geopolitical.languages(language_id),
    currency_id UUID,
    phone_prefix VARCHAR(10) COLLATE "C",
    is_active BOOLEAN DEFAULT true NOT NULL,
    is_deleted BOOLEAN DEFAULT false NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    created_by UUID,
    created_ip INET,
    created_device JSONB,
    created_session UUID,
    created_location JSONB,
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by UUID,
    updated_ip INET,
    updated_device JSONB,
    updated_session UUID,
    updated_location JSONB,
    deleted_at TIMESTAMPTZ,
    deleted_by UUID,
    deleted_ip INET,
    deleted_device JSONB,
    deleted_session UUID,
    deleted_location JSONB,
    source_system VARCHAR(50) DEFAULT 'reference_master_geopolitical',
    change_reason TEXT,
    version INTEGER DEFAULT 1 NOT NULL
);