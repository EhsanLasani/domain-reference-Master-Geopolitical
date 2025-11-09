@echo off
echo ========================================
echo  DATABASE SCHEMA INITIALIZATION
echo ========================================
echo.
echo Database: geopolitical
echo Schema: domain_reference_master_geopolitical
echo.
echo Initializing database schema...
echo.

REM Initialize schema
psql -h localhost -U postgres -d geopolitical -f database-layer\postgresql-primary\00-init-schema.sql

echo.
echo Creating entities...
echo.

REM Create entities in order
psql -h localhost -U postgres -d geopolitical -f database-layer\postgresql-primary\entities\00-schema-setup.sql
psql -h localhost -U postgres -d geopolitical -f database-layer\postgresql-primary\entities\01-regions.sql
psql -h localhost -U postgres -d geopolitical -f database-layer\postgresql-primary\entities\02-languages.sql
psql -h localhost -U postgres -d geopolitical -f database-layer\postgresql-primary\entities\03-timezones.sql
psql -h localhost -U postgres -d geopolitical -f database-layer\postgresql-primary\entities\04-countries.sql
psql -h localhost -U postgres -d geopolitical -f database-layer\postgresql-primary\entities\05-subdivisions.sql
psql -h localhost -U postgres -d geopolitical -f database-layer\postgresql-primary\entities\06-locales.sql

echo.
echo ✅ Database initialization completed!
echo.
echo You can now run the aligned application:
echo   ./run-aligned.bat
echo.
pause