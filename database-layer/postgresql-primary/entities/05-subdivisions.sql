-- ============================================================================
-- ENTITY: country_subdivisions
-- PURPOSE: Country subdivisions table with LASANI compliance
-- DEPENDENCIES: countries
-- ============================================================================

CREATE TABLE domain_reference_master_geopolitical.country_subdivisions (
    subdivision_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subdivision_code VARCHAR(10) COLLATE "C" NOT NULL,
    subdivision_name VARCHAR(100) NOT NULL,
    country_id UUID NOT NULL REFERENCES domain_reference_master_geopolitical.countries(country_id),
    subdivision_type VARCHAR(20),
    parent_subdivision_id UUID REFERENCES domain_reference_master_geopolitical.country_subdivisions(subdivision_id),
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
    version INTEGER DEFAULT 1 NOT NULL,
    UNIQUE(country_id, subdivision_code)
);