@echo off
echo Starting Phase 2 with clean environment...

REM Clear persistent environment variables
set DB_NAME=
set DB_PASSWORD=
set DB_USER=
set DB_HOST=

REM Set correct values
set DB_NAME=referencemaster
set DB_PASSWORD=@Salman2021
set DB_USER=postgres
set DB_HOST=localhost

echo Downloading dependencies...
go mod tidy

echo Starting Phase 2 application with clean environment...
go run cmd/server/main.go