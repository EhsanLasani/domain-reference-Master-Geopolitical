@echo off
echo Starting LASANI Platform - Reference Master Geopolitical Service...
echo.

REM Set environment variables
set DB_HOST=localhost
set DB_PORT=5432
set DB_USER=postgres
set DB_PASSWORD=@Salman2021
set DB_NAME=referencemaster
set DB_SSL_MODE=disable
set REDIS_HOST=localhost
set REDIS_PORT=6379
set SERVER_PORT=8081

echo Environment Configuration:
echo - Database: %DB_HOST%:%DB_PORT%/%DB_NAME%
echo - Redis: %REDIS_HOST%:%REDIS_PORT%
echo - Server Port: %SERVER_PORT%
echo.

echo Building and running Go application...
go mod tidy
go run cmd/server/main.go

pause