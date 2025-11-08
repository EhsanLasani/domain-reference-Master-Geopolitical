@echo off
echo Starting Geopolitical Service on port 8081...

REM Set environment variables directly
set JWT_SECRET=geopolitical-service-jwt-secret-key-2024
set DB_PASSWORD=postgres123
set DB_HOST=localhost
set DB_PORT=5432
set DB_NAME=geopolitical
set DB_USER=postgres
set SERVER_PORT=8081

echo Starting application on port 8081...
go run cmd/server/main.go