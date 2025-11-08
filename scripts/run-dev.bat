@echo off
echo Starting Geopolitical Service in development mode...

REM Check if .env exists
if not exist .env (
    echo Creating .env from template...
    copy .env.example .env
)

REM Start the application
echo Starting application...
go run cmd/server/main.go