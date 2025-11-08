@echo off
echo Starting Complete Geopolitical Reference Service...
echo.

echo Checking database connection...
go run test-db.go
if %ERRORLEVEL% neq 0 (
    echo Database connection failed!
    pause
    exit /b 1
)

echo.
echo Starting application with all entities (Countries, Regions, Languages)...
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=postgres
set DB_PASSWORD=@Salman2021
set DB_NAME=referencemaster
set DB_SSL_MODE=disable
set REDIS_HOST=localhost
set REDIS_PORT=6379
set SERVER_PORT=8081

go run cmd/server/main.go
pause