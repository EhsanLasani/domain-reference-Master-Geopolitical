@echo off
echo Setting up database for Geopolitical Domain Service...

REM Check if PostgreSQL is running
pg_isready -h localhost -p 5432
if %errorlevel% neq 0 (
    echo ERROR: PostgreSQL is not running on localhost:5432
    echo Please start PostgreSQL service first
    pause
    exit /b 1
)

REM Create database if it doesn't exist
psql -h localhost -U postgres -c "CREATE DATABASE referencemaster;" 2>nul
if %errorlevel% equ 0 (
    echo Database 'referencemaster' created successfully
) else (
    echo Database 'referencemaster' already exists or creation failed
)

REM Initialize schema
echo Initializing schema...
psql -h localhost -U postgres -d referencemaster -f "database-layer\postgresql-primary\00-init-schema.sql"

REM Create entities
echo Creating entities...
for %%f in (database-layer\postgresql-primary\entities\*.sql) do (
    echo Processing %%f...
    psql -h localhost -U postgres -d referencemaster -f "%%f"
)

REM Create indexes
echo Creating indexes...
for %%f in (database-layer\postgresql-primary\indexes\*.sql) do (
    echo Processing %%f...
    psql -h localhost -U postgres -d referencemaster -f "%%f"
)

REM Seed initial data
echo Seeding initial data...
psql -h localhost -U postgres -d referencemaster -f "database-layer\postgresql-primary\seed-data.sql"

echo Database setup completed successfully!
echo You can now run: go run cmd\server\main.go
pause