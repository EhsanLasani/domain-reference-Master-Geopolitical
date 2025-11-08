# Schema Changes Summary - Countries Table Updates

## Changes Made

### 1. **Foreign Key Relationships**
- **region_code** → **region_id UUID** (FK to regions.region_id)
- **currency_code** → **currency_id UUID** (FK to commerce domain)
- **Added**: **primary_language_id UUID** (FK to languages.language_id)

### 2. **Continent Code Constraint**
- **continent_code** now has ENUM constraint: `('AF', 'AS', 'EU', 'NA', 'SA', 'OC', 'AN')`

### 3. **Updated Indexes**
- Removed: `idx_countries_currency`, `idx_countries_region`
- Added: `idx_countries_continent`, `idx_countries_region_fk`, `idx_countries_language_fk`, `idx_countries_currency_fk`

### 4. **Updated Views**
- **v_countries_with_regions** → **v_countries_with_details** (includes language and proper FK joins)

### 5. **Helper Function Updates**
- **insert_country_dev()**: Updated to accept region_code/language_code and lookup UUIDs
- **bulk_insert_countries_dev()**: Updated CSV format

## New CSV Format for Countries

```csv
country_code,country_name,iso3_code,numeric_code,official_name,capital_city,continent_code,region_code,language_code,currency_id,phone_prefix
US,United States,USA,840,United States of America,Washington D.C.,NA,NORTH_AM,en,uuid-here,+1
```

## Migration Required

### Before Using Updated Schema:
1. **Deploy regions and languages first**
2. **Update existing countries data** to use UUIDs instead of codes
3. **Commerce domain** must provide currency_id values

### Continent Code Values:
- **AF**: Africa
- **AS**: Asia  
- **EU**: Europe
- **NA**: North America
- **SA**: South America
- **OC**: Oceania
- **AN**: Antarctica

## Impact on Development Import

### Helper Functions Still Work:
```sql
-- Lookup by codes (recommended)
SELECT insert_country_dev('US', 'United States', 'USA', 840, 'United States of America', 'Washington D.C.', 'NA', 'NORTH_AM', 'en', null, '+1');

-- Bulk import with code lookup
SELECT bulk_insert_countries_dev('
US,United States,USA,840,United States of America,Washington D.C.,NA,NORTH_AM,en,,+1
CA,Canada,CAN,124,Canada,Ottawa,NA,NORTH_AM,en,,+1
');
```

### Direct SQL Requires UUIDs:
```sql
-- Must use actual UUIDs for foreign keys
INSERT INTO countries (country_code, country_name, region_id, primary_language_id) 
VALUES ('US', 'United States', 'uuid-of-north-america-region', 'uuid-of-english-language');
```

## Status
- ✅ Schema updated with proper foreign keys
- ✅ Helper functions updated with code lookup
- ✅ Views updated for new relationships
- ⚠️ **Requires regions and languages data before countries**
- ⚠️ **Currency_id will be NULL until commerce domain is implemented**