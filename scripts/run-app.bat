@echo off
echo Starting Geopolitical Service...

REM Set environment variables directly
set JWT_SECRET=geopolitical-service-jwt-secret-key-2024
set DB_PASSWORD=postgres123
set DB_HOST=localhost
set DB_PORT=5432
set DB_NAME=geopolitical
set DB_USER=postgres

echo Starting application with environment variables...
go run cmd/server/main.go