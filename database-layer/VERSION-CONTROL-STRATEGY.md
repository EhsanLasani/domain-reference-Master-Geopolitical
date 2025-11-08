# Database Version Control Strategy

## Current Situation
You ran the **first version** of `schema-definition.sql` in DBeaver. The current file has been updated with additional features.

## Solution: Migration-Based Approach

### 1. **Current Database State**
```sql
-- What you have deployed:
- Tables: 6 (countries, subdivisions, regions, timezones, languages, locales)
- Basic indexes: Business identity + foreign keys
- LASANI audit fields: ✅ Complete
- Missing: LASANI audit indexes, constraints, views
```

### 2. **Apply Updates**
Run the new migration file in DBeaver:
```bash
# File to run: 03-schema-updates.sql
# This adds missing components without breaking existing data
```

## Database Migration Strategy

### File Organization
```
database-layer/postgresql-primary/
├── schema-definition.sql          # Complete schema (for new deployments)
├── 02-seed-data.sql              # Initial data
├── 03-schema-updates.sql         # Updates for existing deployments
└── migrations/
    ├── 001_initial_schema.sql    # Future: Initial schema
    ├── 002_add_constraints.sql   # Future: Add constraints
    ├── 003_add_audit_indexes.sql # Future: Add audit indexes
    └── 004_add_views.sql         # Future: Add views
```

### Version Control Fields

#### Table-Level Versioning
```sql
-- Each record has version field for optimistic locking
version INTEGER DEFAULT 1 NOT NULL

-- Usage in application:
UPDATE countries 
SET country_name = 'New Name', version = version + 1
WHERE country_id = ? AND version = ?  -- Prevents concurrent updates
```

#### Schema-Level Versioning
```sql
-- Create schema version tracking table
CREATE TABLE IF NOT EXISTS domain_reference_master_geopolitical.schema_versions (
    version_id SERIAL PRIMARY KEY,
    version_number VARCHAR(20) NOT NULL,
    description TEXT,
    applied_at TIMESTAMPTZ DEFAULT NOW(),
    applied_by VARCHAR(100),
    script_name VARCHAR(255)
);

-- Track current deployment
INSERT INTO domain_reference_master_geopolitical.schema_versions 
(version_number, description, script_name)
VALUES 
('1.0.0', 'Initial schema deployment', 'schema-definition.sql'),
('1.1.0', 'Added constraints, indexes, and views', '03-schema-updates.sql');
```

## Migration Execution Plan

### Step 1: Apply Current Updates
```sql
-- Run in DBeaver:
\i 03-schema-updates.sql
```

### Step 2: Verify Updates
```sql
-- Check indexes were created
SELECT schemaname, tablename, indexname 
FROM pg_indexes 
WHERE schemaname = 'domain_reference_master_geopolitical'
ORDER BY tablename, indexname;

-- Check constraints were added
SELECT conname, conrelid::regclass 
FROM pg_constraint 
WHERE connamespace = 'domain_reference_master_geopolitical'::regnamespace;

-- Check views were created
SELECT viewname 
FROM pg_views 
WHERE schemaname = 'domain_reference_master_geopolitical';
```

### Step 3: Future Migrations
```sql
-- For future changes, create numbered migration files:
-- 004_add_new_table.sql
-- 005_modify_constraints.sql
-- etc.
```

## Version Control Best Practices

### 1. **Never Modify Deployed Files**
- ❌ Don't change `schema-definition.sql` after deployment
- ✅ Create new migration files for changes

### 2. **Always Use IF NOT EXISTS**
```sql
-- Safe for re-running
CREATE INDEX IF NOT EXISTS idx_name ON table(column);
ALTER TABLE table ADD CONSTRAINT IF NOT EXISTS chk_name CHECK (condition);
CREATE OR REPLACE VIEW view_name AS SELECT ...;
```

### 3. **Track All Changes**
```sql
-- Log every migration
INSERT INTO schema_versions (version_number, description, script_name)
VALUES ('1.2.0', 'Added new business rules', '004_business_rules.sql');
```

### 4. **Rollback Strategy**
```sql
-- Each migration should have rollback script
-- 004_business_rules.sql -> 004_business_rules_rollback.sql
```

## Application-Level Version Control

### Optimistic Locking Pattern
```go
// Go application code
func UpdateCountry(ctx context.Context, country *Country) error {
    result := db.Model(country).
        Where("country_id = ? AND version = ?", country.ID, country.Version).
        Updates(map[string]interface{}{
            "country_name": country.Name,
            "updated_at":   time.Now(),
            "updated_by":   getCurrentUser(ctx),
            "version":      gorm.Expr("version + 1"),
        })
    
    if result.RowsAffected == 0 {
        return errors.New("version conflict - record was modified by another user")
    }
    
    country.Version++ // Update local version
    return nil
}
```

### Audit Trail Integration
```go
// Automatic audit field population
func (c *Country) BeforeUpdate(tx *gorm.DB) error {
    c.UpdatedAt = time.Now()
    c.UpdatedBy = getCurrentUserID(tx.Statement.Context)
    c.UpdatedIP = getCurrentIP(tx.Statement.Context)
    // ... other audit fields
    return nil
}
```

## Current Action Required

### Immediate Steps
1. **Run `03-schema-updates.sql` in DBeaver**
2. **Verify all components are added**
3. **Test basic queries on views**

### Commands to Execute
```sql
-- 1. Apply updates
\i 03-schema-updates.sql

-- 2. Verify deployment
SELECT 'Indexes' as component, COUNT(*) as count 
FROM pg_indexes 
WHERE schemaname = 'domain_reference_master_geopolitical'
UNION ALL
SELECT 'Constraints', COUNT(*) 
FROM pg_constraint 
WHERE connamespace = 'domain_reference_master_geopolitical'::regnamespace
UNION ALL
SELECT 'Views', COUNT(*) 
FROM pg_views 
WHERE schemaname = 'domain_reference_master_geopolitical';

-- 3. Test views
SELECT * FROM domain_reference_master_geopolitical.v_audit_summary;
```

**This approach ensures your existing data is preserved while adding all the missing enterprise features.**