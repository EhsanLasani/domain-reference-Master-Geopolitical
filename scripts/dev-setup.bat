@echo off
echo Setting up development environment...

REM Copy environment file if it doesn't exist
if not exist .env (
    copy .env.example .env
    echo Created .env file from template
)

REM Download Go dependencies
echo Downloading Go dependencies...
go mod download
go mod tidy

echo Development environment ready!
echo.
echo Next steps:
echo 1. Update .env file with your configuration
echo 2. Run: docker-compose up -d
echo 3. Run: go run cmd/server/main.go