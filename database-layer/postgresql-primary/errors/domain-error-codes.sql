-- ============================================================================
-- DOMAIN ERROR CODES
-- PURPOSE: Standardized error codes for geopolitical domain
-- ============================================================================

CREATE TABLE IF NOT EXISTS domain_reference_master_geopolitical.error_codes (
    error_code VARCHAR(10) PRIMARY KEY,
    error_category VARCHAR(50) NOT NULL,
    error_message TEXT NOT NULL,
    http_status_code INTEGER NOT NULL,
    is_retryable BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now()
);

-- Geopolitical domain error codes (GEO-1xxx)
INSERT INTO domain_reference_master_geopolitical.error_codes VALUES
('GEO-1001', 'VALIDATION', 'Invalid country code format', 400, false, now()),
('GEO-1002', 'VALIDATION', 'Country code already exists', 409, false, now()),
('GEO-1003', 'VALIDATION', 'Invalid continent code', 400, false, now()),
('GEO-1004', 'BUSINESS', 'Cannot delete active country', 400, false, now()),
('GEO-1005', 'NOT_FOUND', 'Country not found', 404, false, now()),
('GEO-1006', 'DATABASE', 'Database connection failed', 503, true, now()),
('GEO-1007', 'DATABASE', 'Query timeout', 504, true, now()),
('GEO-1008', 'DATABASE', 'Constraint violation', 400, false, now()),
('GEO-1009', 'BUSINESS', 'Version conflict', 409, true, now()),
('GEO-1010', 'AUTHORIZATION', 'Insufficient permissions', 403, false, now())
ON CONFLICT (error_code) DO NOTHING;