-- ============================================================================
-- MIGRATION SYSTEM
-- PURPOSE: Track schema migrations and versions
-- ============================================================================

CREATE TABLE IF NOT EXISTS domain_reference_master_geopolitical.schema_migrations (
    version VARCHAR(50) PRIMARY KEY,
    description TEXT NOT NULL,
    applied_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    applied_by VARCHAR(100) DEFAULT current_user NOT NULL,
    checksum VARCHAR(64),
    execution_time_ms INTEGER,
    rollback_sql TEXT
);

CREATE TABLE IF NOT EXISTS domain_reference_master_geopolitical.query_performance_log (
    log_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sql_key VARCHAR(100) NOT NULL,
    operation VARCHAR(20) NOT NULL,
    duration_ms INTEGER NOT NULL,
    tenant_id UUID,
    executed_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    query_hash VARCHAR(64),
    result_count INTEGER,
    error_message TEXT
);

CREATE INDEX IF NOT EXISTS idx_query_performance_sql_key 
    ON domain_reference_master_geopolitical.query_performance_log(sql_key);
CREATE INDEX IF NOT EXISTS idx_query_performance_tenant 
    ON domain_reference_master_geopolitical.query_performance_log(tenant_id);
CREATE INDEX IF NOT EXISTS idx_query_performance_executed_at 
    ON domain_reference_master_geopolitical.query_performance_log(executed_at);

-- Insert initial migration record
INSERT INTO domain_reference_master_geopolitical.schema_migrations 
(version, description, checksum) 
VALUES 
('001', 'Initial schema migrations table', 'abc123')
ON CONFLICT (version) DO NOTHING;