# Development Data Import Workflow

## Overview
Complete workflow for importing reference data during development phase with automatic LASANI audit compliance.

## How to Use This Workflow

### For New Developers
1. **Read this document completely** before starting
2. **Follow Initial Setup** (one-time only)
3. **Set your developer context** (required)
4. **Choose import method** based on your data source
5. **Verify data quality** after import

### For Existing Developers
1. **Set/switch context** for your current work
2. **Import data** using any method
3. **Clean up** when switching tasks

### Document Sections
- **Initial Setup**: First-time database configuration
- **Daily Workflow**: Regular data import operations
- **Context Management**: Developer isolation and switching
- **Adding New Tables**: Schema expansion workflow
- **Troubleshooting**: Common issues and fixes
- **Quick Reference**: Command cheat sheet

## Prerequisites
- PostgreSQL 15+ with pgcrypto extension
- DBeaver or similar database tool
- Access to domain schema

## Initial Setup (One-Time)

### 1. Deploy Database Schema
```bash
# Navigate to database layer
cd domain-reference-Master-Geopolitical/database-layer/postgresql-primary/

# Execute files in order
psql -h localhost -U postgres -d lasani -f schema-definition.sql
psql -h localhost -U postgres -d lasani -f 02-development-data-import-helpers.sql
psql -h localhost -U postgres -d lasani -f 04-auto-generate-helpers.sql
psql -h localhost -U postgres -d lasani -f 05-schema-change-detection.sql
```

### 2. Set Developer Context
```sql
-- REQUIRED: Set your developer identity (replace with your details)
SELECT set_dev_context('john_doe', 'john.doe@company.com', 'main');
```

### 3. Verify Setup
```sql
-- Check context is active
SELECT * FROM list_dev_contexts();

-- Verify triggers are installed
SELECT * FROM detect_and_generate_helpers();
```

## Daily Development Workflow

### Option A: Helper Functions (Recommended)

#### Single Record Insert
```sql
-- Basic insert (only required fields)
SELECT insert_country_dev('US', 'United States');

-- Full insert (all fields)
SELECT insert_country_dev('US', 'United States', 'USA', 840, 'United States of America', 'Washington D.C.', 'NA', 'NORTH_AM', 'USD', '+1');
```

#### Bulk CSV Import
```sql
-- Method 1: Direct CSV string
SELECT bulk_insert_countries_dev('
country_code,country_name,iso3_code,numeric_code
US,United States,USA,840
CA,Canada,CAN,124
GB,United Kingdom,GBR,826
');

-- Method 2: Minimal format (only required fields)
SELECT bulk_insert_countries_dev('
US,United States
CA,Canada
GB,United Kingdom
');

-- Method 3: Using DO block for large datasets
DO $$
DECLARE
    csv_data TEXT := '
country_code,country_name,iso3_code,numeric_code
SG,Singapore,SGP,702
MY,Malaysia,MYS,458
TH,Thailand,THA,764
';
    inserted_count INTEGER;
BEGIN
    SELECT bulk_insert_countries_dev(csv_data) INTO inserted_count;
    RAISE NOTICE 'Inserted % countries', inserted_count;
END $$;
```

#### All Available Bulk Functions
```sql
-- Countries
SELECT bulk_insert_countries_dev('code,name,iso3,numeric\nUS,United States,USA,840');

-- Regions  
SELECT bulk_insert_regions_dev('code,name,type\nNA,North America,Continent');

-- Languages
SELECT bulk_insert_languages_dev('code,name,iso3,native,direction\nen,English,eng,English,LTR');

-- Timezones
SELECT bulk_insert_timezones_dev('code,name,hours,minutes,dst,dst_hours\nAmerica/New_York,Eastern Time,-5,0,true,-4');

-- Subdivisions (requires countries first)
SELECT bulk_insert_subdivisions_dev('code,name,country_code,type\nCA,California,US,State');

-- Locales (requires languages and countries first)
SELECT bulk_insert_locales_dev('code,name,language_code,country_code\nen-US,English (US),en,US');
```

### Option B: Direct SQL (Auto-Audit Enabled)
```sql
-- Direct INSERT - audit fields auto-populated by triggers
INSERT INTO domain_reference_master_geopolitical.countries (country_code, country_name) 
VALUES ('XX', 'Test Country');

-- Bulk INSERT - also works with auto-audit
INSERT INTO domain_reference_master_geopolitical.countries (country_code, country_name) 
VALUES 
('YY', 'Country Y'),
('ZZ', 'Country Z');
```

### Option C: DBeaver Import Wizard

#### Step-by-Step Process
1. **Right-click table** → "Import Data"
2. **Select CSV file** or paste data
3. **Map columns**:
   - Map business columns (country_code, country_name, etc.)
   - **SKIP audit fields** (created_by, created_at, etc.) - auto-populated
   - **SKIP status fields** (is_active, is_deleted) - have defaults
4. **Preview and execute**

#### CSV File Format Example
```csv
country_code,country_name,iso3_code,numeric_code,official_name,capital_city
US,United States,USA,840,United States of America,Washington D.C.
CA,Canada,CAN,124,Canada,Ottawa
GB,United Kingdom,GBR,826,United Kingdom of Great Britain and Northern Ireland,London
```

#### DBeaver Mapping Tips
- **Required fields**: country_code, country_name (minimum)
- **Optional fields**: All others can be left unmapped
- **Auto-populated**: All LASANI audit fields (27 fields)
- **Defaults**: is_active=true, is_deleted=false, version=1

## Context Management

### Switch Between Contexts
```sql
-- Switch to feature branch
SELECT switch_dev_context('john_doe', 'feature_xyz');

-- Switch back to main
SELECT switch_dev_context('john_doe', 'main');
```

### View Your Data
```sql
-- Statistics by developer
SELECT * FROM get_dev_data_stats('john_doe');

-- All contexts overview
SELECT * FROM list_dev_contexts();
```

### Clean Up Data
```sql
-- Clear specific context
SELECT clear_dev_data_by_context('john_doe', 'feature_xyz');

-- Clear current active context
SELECT clear_dev_data();
```

## Adding New Tables

### When Schema Changes
```sql
-- After adding new table to schema-definition.sql
SELECT detect_and_generate_helpers();

-- Verify new helpers were created
SELECT * FROM domain_reference_master_geopolitical.schema_tracking;
```

### Manual Helper Generation
```sql
-- Generate for specific table
SELECT generate_dev_helpers_for_table('domain_reference_master_geopolitical', 'new_table_name');

-- Generate for all tables
SELECT generate_all_dev_helpers();
```

## Data Verification

### Check Data Quality
```sql
-- View active records
SELECT * FROM domain_reference_master_geopolitical.v_countries_active LIMIT 10;

-- Audit trail verification
SELECT * FROM domain_reference_master_geopolitical.v_countries_audit_trail LIMIT 5;

-- Domain summary
SELECT * FROM domain_reference_master_geopolitical.v_audit_summary;
```

### Validate LASANI Compliance
```sql
-- Check audit field population
SELECT 
    country_code,
    created_by IS NOT NULL as has_created_by,
    created_at IS NOT NULL as has_created_at,
    created_device IS NOT NULL as has_device_info
FROM domain_reference_master_geopolitical.countries 
LIMIT 5;
```

## Troubleshooting

### Common Issues

#### 1. Audit Fields Not Populated
```sql
-- Check if triggers exist
SELECT trigger_name, event_manipulation, action_statement 
FROM information_schema.triggers 
WHERE trigger_schema = 'domain_reference_master_geopolitical';

-- Recreate triggers if missing
\i 02-development-data-import-helpers.sql
```

#### 2. Context Not Set
```sql
-- Check current context
SELECT get_current_dev_context();

-- Reset to default
SELECT set_dev_context('system', 'dev@lasani.com', 'default');
```

#### 3. Helper Functions Missing
```sql
-- Check what helpers exist
SELECT routine_name 
FROM information_schema.routines 
WHERE routine_schema = 'domain_reference_master_geopolitical' 
AND routine_name LIKE '%_dev';

-- Regenerate all helpers
SELECT generate_all_dev_helpers();
```

## Best Practices

### Development Workflow
1. **Always set developer context first**
2. **Use meaningful context names** (feature branches, testing, etc.)
3. **Clean up feature branch data** when done
4. **Verify data quality** before sharing

### Data Import Guidelines
1. **Use helper functions** for complex data with relationships
2. **Use direct SQL** for simple bulk imports
3. **Use DBeaver wizard** for CSV files
4. **Always verify** audit fields are populated

### Team Collaboration
1. **Use unique developer names** (no conflicts)
2. **Document data sources** in change_reason field
3. **Share context names** for collaborative work
4. **Regular cleanup** of old contexts

## File Structure
```
domain-reference-Master-Geopolitical/database-layer/
├── postgresql-primary/
│   ├── schema-definition.sql                    # Core schema
│   ├── 02-development-data-import-helpers.sql   # Helper functions & triggers
│   ├── 03-sample-seed-data.sql                 # Sample data
│   ├── 04-auto-generate-helpers.sql            # Auto-generation tools
│   └── 05-schema-change-detection.sql          # Schema change detection
└── DEVELOPMENT-IMPORT-WORKFLOW.md              # This document
```

## Quick Reference

### Essential Commands
```sql
-- Setup
SELECT set_dev_context('your_name', 'your@email.com', 'context');

-- Import
SELECT insert_country_dev('CODE', 'Name');
SELECT bulk_insert_countries_dev('CSV_DATA');

-- Manage
SELECT * FROM get_dev_data_stats('your_name');
SELECT clear_dev_data_by_context('your_name', 'context');

-- New Tables
SELECT detect_and_generate_helpers();
```

### All Available Bulk Import Functions
- `bulk_insert_countries_dev(csv_data)`
- `bulk_insert_regions_dev(csv_data)`
- `bulk_insert_languages_dev(csv_data)`
- `bulk_insert_timezones_dev(csv_data)`
- `bulk_insert_subdivisions_dev(csv_data)`
- `bulk_insert_locales_dev(csv_data)`

---

**Status**: Ready for development phase data import with full LASANI audit compliance.