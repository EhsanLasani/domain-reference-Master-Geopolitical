-- ============================================================================
-- ROW LEVEL SECURITY POLICIES
-- PURPOSE: Tenant isolation and security enforcement
-- ============================================================================

-- Enable RLS on all tables
ALTER TABLE domain_reference_master_geopolitical.countries ENABLE ROW LEVEL SECURITY;
ALTER TABLE domain_reference_master_geopolitical.regions ENABLE ROW LEVEL SECURITY;
ALTER TABLE domain_reference_master_geopolitical.languages ENABLE ROW LEVEL SECURITY;
ALTER TABLE domain_reference_master_geopolitical.timezones ENABLE ROW LEVEL SECURITY;
ALTER TABLE domain_reference_master_geopolitical.subdivisions ENABLE ROW LEVEL SECURITY;
ALTER TABLE domain_reference_master_geopolitical.locales ENABLE ROW LEVEL SECURITY;

-- Tenant isolation policy
CREATE POLICY tenant_isolation ON domain_reference_master_geopolitical.countries
    FOR ALL TO PUBLIC
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

CREATE POLICY tenant_isolation ON domain_reference_master_geopolitical.regions
    FOR ALL TO PUBLIC
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

CREATE POLICY tenant_isolation ON domain_reference_master_geopolitical.languages
    FOR ALL TO PUBLIC
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

CREATE POLICY tenant_isolation ON domain_reference_master_geopolitical.timezones
    FOR ALL TO PUBLIC
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

CREATE POLICY tenant_isolation ON domain_reference_master_geopolitical.subdivisions
    FOR ALL TO PUBLIC
    USING (tenant_id = current_setting('app.tenant_id')::uuid);

CREATE POLICY tenant_isolation ON domain_reference_master_geopolitical.locales
    FOR ALL TO PUBLIC
    USING (tenant_id = current_setting('app.tenant_id')::uuid);