-- ============================================================================
-- ENTITY: timezones
-- PURPOSE: Timezones table with LASANI compliance
-- DEPENDENCIES: None (foundational entity)
-- ============================================================================

CREATE TABLE domain_reference_master_geopolitical.timezones (
    timezone_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timezone_code VARCHAR(50) COLLATE "C" UNIQUE NOT NULL,
    timezone_name VARCHAR(100) NOT NULL,
    utc_offset_hours SMALLINT NOT NULL,
    utc_offset_minutes SMALLINT DEFAULT 0,
    supports_dst BOOLEAN DEFAULT false,
    dst_offset_hours SMALLINT,
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