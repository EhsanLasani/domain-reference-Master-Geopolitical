-- ============================================================================
-- TENANT MIGRATION SCRIPT
-- PURPOSE: Add tenant_id column to all existing tables
-- SCHEMA: domain_reference_master_geopolitical
-- ============================================================================

-- Set search path
SET search_path TO domain_reference_master_geopolitical, public;

-- Add tenant_id to countries table if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'domain_reference_master_geopolitical' 
        AND table_name = 'countries' 
        AND column_name = 'tenant_id'
    ) THEN
        ALTER TABLE domain_reference_master_geopolitical.countries 
        ADD COLUMN tenant_id VARCHAR(100) DEFAULT 'default-tenant';
        
        -- Update existing records
        UPDATE domain_reference_master_geopolitical.countries 
        SET tenant_id = 'default-tenant' 
        WHERE tenant_id IS NULL OR tenant_id = '';
        
        -- Make column NOT NULL
        ALTER TABLE domain_reference_master_geopolitical.countries 
        ALTER COLUMN tenant_id SET NOT NULL;
        
        -- Add index
        CREATE INDEX IF NOT EXISTS idx_countries_tenant_id 
        ON domain_reference_master_geopolitical.countries(tenant_id);
        
        RAISE NOTICE 'Added tenant_id to countries table';
    END IF;
END $$;

-- Add tenant_id to regions table if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'domain_reference_master_geopolitical' 
        AND table_name = 'regions' 
        AND column_name = 'tenant_id'
    ) THEN
        ALTER TABLE domain_reference_master_geopolitical.regions 
        ADD COLUMN tenant_id VARCHAR(100) DEFAULT 'default-tenant';
        
        UPDATE domain_reference_master_geopolitical.regions 
        SET tenant_id = 'default-tenant' 
        WHERE tenant_id IS NULL OR tenant_id = '';
        
        ALTER TABLE domain_reference_master_geopolitical.regions 
        ALTER COLUMN tenant_id SET NOT NULL;
        
        CREATE INDEX IF NOT EXISTS idx_regions_tenant_id 
        ON domain_reference_master_geopolitical.regions(tenant_id);
        
        RAISE NOTICE 'Added tenant_id to regions table';
    END IF;
END $$;

-- Add tenant_id to languages table if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'domain_reference_master_geopolitical' 
        AND table_name = 'languages' 
        AND column_name = 'tenant_id'
    ) THEN
        ALTER TABLE domain_reference_master_geopolitical.languages 
        ADD COLUMN tenant_id VARCHAR(100) DEFAULT 'default-tenant';
        
        UPDATE domain_reference_master_geopolitical.languages 
        SET tenant_id = 'default-tenant' 
        WHERE tenant_id IS NULL OR tenant_id = '';
        
        ALTER TABLE domain_reference_master_geopolitical.languages 
        ALTER COLUMN tenant_id SET NOT NULL;
        
        CREATE INDEX IF NOT EXISTS idx_languages_tenant_id 
        ON domain_reference_master_geopolitical.languages(tenant_id);
        
        RAISE NOTICE 'Added tenant_id to languages table';
    END IF;
END $$;

-- Add tenant_id to timezones table if it doesn't exist
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_schema = 'domain_reference_master_geopolitical' 
        AND table_name = 'timezones' 
        AND column_name = 'tenant_id'
    ) THEN
        ALTER TABLE domain_reference_master_geopolitical.timezones 
        ADD COLUMN tenant_id VARCHAR(100) DEFAULT 'default-tenant';
        
        UPDATE domain_reference_master_geopolitical.timezones 
        SET tenant_id = 'default-tenant' 
        WHERE tenant_id IS NULL OR tenant_id = '';
        
        ALTER TABLE domain_reference_master_geopolitical.timezones 
        ALTER COLUMN tenant_id SET NOT NULL;
        
        CREATE INDEX IF NOT EXISTS idx_timezones_tenant_id 
        ON domain_reference_master_geopolitical.timezones(tenant_id);
        
        RAISE NOTICE 'Added tenant_id to timezones table';
    END IF;
END $$;

-- Verify tenant migration
SELECT 
    table_name,
    column_name,
    is_nullable,
    column_default
FROM information_schema.columns 
WHERE table_schema = 'domain_reference_master_geopolitical' 
AND column_name = 'tenant_id'
ORDER BY table_name;

SELECT 'Tenant migration completed successfully' as status;