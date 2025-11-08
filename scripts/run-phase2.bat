@echo off
echo Starting Phase 2: Complete Enterprise Stack...

REM Set all environment variables
set JWT_SECRET=geopolitical-service-jwt-secret-key-2024
set DB_PASSWORD=@Salman2021
set DB_HOST=localhost
set DB_PORT=5432
set DB_NAME=referencemaster
set DB_USER=postgres
set DB_SSL_MODE=disable
set DB_MAX_CONNECTIONS=25
set REDIS_HOST=localhost
set REDIS_PORT=6379
set REDIS_PASSWORD=
set REDIS_DB=0
set SERVER_PORT=8081
set TRACING_ENABLED=true
set SERVICE_NAME=geopolitical-service

echo Downloading dependencies...
go mod tidy

echo Starting Phase 2 application with full stack...
go run cmd/server/main.go