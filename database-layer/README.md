# Reference Master Geopolitical - Database Layer

## Overview
Complete database implementation for global geopolitical and localization reference data with full LASANI audit compliance.

## Schema Information
- **Schema Name**: `domain_reference_master_geopolitical`
- **Domain**: Reference Master Geopolitical
- **Dependencies**: None (foundational domain)
- **Source System**: `reference_master_geopolitical`

## Tables (6 Total)

### Core Entities
1. **countries** - ISO 3166-1 country reference data
2. **country_subdivisions** - States, provinces, regions within countries
3. **regions** - Continents, economic zones, political regions
4. **timezones** - IANA timezone reference data
5. **languages** - ISO 639 language reference data
6. **locales** - Language-country combinations (RFC 5646)

## LASANI Compliance ✅

### Field Structure (35 fields per table)
```yaml
field_breakdown:
  primary_key: 1        # {entity}_id UUID
  business_identity: 2  # {entity}_code, {entity}_name
  business_properties: 8 # Domain-specific fields (varies by table)
  status_management: 2  # is_active, is_deleted
  lasani_audit: 21     # Complete audit trail
  version_control: 1   # version INTEGER
  total: 35            # 13 business + 22 LASANI
```

### Audit Trail Features
- **Creation Tracking**: User, timestamp, IP, device, session, location
- **Update Tracking**: User, timestamp, IP, device, session, location  
- **Deletion Tracking**: User, timestamp, IP, device, session, location
- **System Metadata**: Source system, change reason
- **Version Control**: Optimistic locking support

## International Compliance ✅

### Collation Strategy
```sql
-- Business codes (case-sensitive, exact matching)
COLLATE "C"

-- Display names (case-insensitive, user-friendly)
COLLATE "en_US.utf8"

-- Native content (locale-aware, international)
COLLATE "und-x-icu"
```

## Performance Optimization ✅

### Index Coverage (84 indexes total)
- **Business Identity**: 12 indexes (code, name per table)
- **Foreign Keys**: 6 indexes (all FK relationships)
- **Status Management**: 6 indexes (active/deleted filtering)
- **LASANI Audit**: 42 indexes (temporal, user, system tracking)
- **Domain-Specific**: 18 indexes (business query patterns)

### Query Views (7 views)
- `v_countries_active` - Active countries only
- `v_countries_with_regions` - Countries with region info
- `v_subdivisions_active` - Active subdivisions with country
- `v_locales_active` - Active locales with language/country
- `v_timezones_active` - Active timezones with formatted offset
- `v_countries_audit_trail` - Compliance audit trail
- `v_audit_summary` - Cross-entity audit summary

## Business Constraints ✅

### Format Validation
```sql
-- Country codes: ISO 3166-1 alpha-2 (US, CA, GB)
country_code ~ '^[A-Z]{2}$'

-- Language codes: ISO 639-1 (en, fr, es)
language_code ~ '^[a-z]{2}$'

-- Locale codes: RFC 5646 (en-US, fr-CA)
locale_code ~ '^[a-z]{2}(-[A-Z]{2})?$'

-- Phone prefixes: International format (+1, +44)
phone_prefix ~ '^\\+?[0-9]{1,4}$'
```

### Business Logic
```sql
-- No self-reference in hierarchies
subdivision_id != parent_subdivision_id
region_id != parent_region_id

-- Valid timezone offsets
utc_offset_hours BETWEEN -12 AND 14
utc_offset_minutes IN (0, 15, 30, 45)

-- Valid text directions
direction IN ('LTR', 'RTL')
```

### Lifecycle Management
```sql
-- Cannot be active and deleted simultaneously
NOT (is_active = true AND is_deleted = true)
```

## Data Access Patterns

### Common Queries
```sql
-- Active countries by region
SELECT * FROM v_countries_active WHERE region_code = 'EU';

-- Country subdivisions
SELECT * FROM v_subdivisions_active WHERE country_code = 'US';

-- Timezone lookup
SELECT * FROM v_timezones_active WHERE timezone_code = 'America/New_York';

-- Locale information
SELECT * FROM v_locales_active WHERE locale_code = 'en-US';
```

### Audit Queries
```sql
-- Recent changes
SELECT * FROM v_countries_audit_trail WHERE updated_at > NOW() - INTERVAL '7 days';

-- Entity statistics
SELECT * FROM v_audit_summary ORDER BY total_records DESC;
```

## Integration Points

### Cross-Domain References
```sql
-- Used by other domains
domain_reference_master_commerce.currencies(country_id)
domain_user_management.users(country_id, language_id, timezone_id)
domain_tenant_management.tenants(country_id, locale_id)
```

### API Integration
- **GraphQL**: Primary interface (80% of queries)
- **REST**: Bulk operations and exports (20% of queries)
- **Repository Pattern**: Data access abstraction

## Deployment

### Schema Deployment
```bash
# Deploy complete schema
psql -h localhost -U postgres -d postgres -f schema-definition.sql

# Load seed data
psql -h localhost -U postgres -d postgres -f 02-seed-data.sql
```

### Index Verification
```sql
-- Verify index creation
\di domain_reference_master_geopolitical.*

-- Check index usage
SELECT schemaname, tablename, indexname, idx_scan 
FROM pg_stat_user_indexes 
WHERE schemaname = 'domain_reference_master_geopolitical'
ORDER BY idx_scan DESC;
```

## Quality Assurance

### Compliance Verification
- ✅ LASANI 27-field audit system
- ✅ International collation support
- ✅ Complete indexing strategy
- ✅ Business constraint validation
- ✅ Audit trail views
- ✅ Performance optimization

### Testing Coverage
- Schema structure validation
- Constraint enforcement testing
- Index performance verification
- Collation behavior validation
- Audit trail functionality
- Cross-domain integration

## Maintenance

### Regular Tasks
- Monitor index usage statistics
- Verify constraint enforcement
- Audit trail compliance checks
- Performance optimization reviews
- Schema evolution management

### Monitoring Queries
```sql
-- Table sizes
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) 
FROM pg_tables WHERE schemaname = 'domain_reference_master_geopolitical';

-- Constraint violations
SELECT conname, conrelid::regclass FROM pg_constraint 
WHERE connamespace = 'domain_reference_master_geopolitical'::regnamespace;
```

## Next Steps
1. Deploy to development environment
2. Load production seed data
3. Implement data access layer (auto-generated models/repositories)
4. Create business logic layer (domain services)
5. Build presentation layer (GraphQL/REST APIs)

**This database layer provides the foundation for all geopolitical reference data across the enterprise platform.**