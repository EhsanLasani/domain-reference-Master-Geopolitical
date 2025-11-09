-- ============================================================================
-- SCHEMA ORCHESTRATOR: Reference Master Geopolitical
-- PURPOSE: Coordinates entity creation following SRP principles
-- DEPENDENCIES: Individual entity files in correct order
-- ============================================================================

-- Entity creation in dependency order
\i entities/00-schema-setup.sql
\i entities/01-regions.sql
\i entities/02-languages.sql
\i entities/03-timezones.sql
\i entities/04-countries.sql
\i entities/05-subdivisions.sql
\i entities/06-locales.sql

-- Index creation per entity
\i indexes/regions-indexes.sql
\i indexes/countries-indexes.sql

-- Security policies
\i security/rls-policies.sql

-- Error codes
\i errors/domain-error-codes.sql

-- Migration system
\i migrations/001_schema_migrations.sql