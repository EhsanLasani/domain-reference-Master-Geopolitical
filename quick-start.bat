@echo off
echo ========================================
echo LASANI Geopolitical Domain Quick Start
echo ========================================

REM Step 1: Check Prerequisites
echo [1/6] Checking prerequisites...
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go 1.21+ from https://golang.org/dl/
    pause
    exit /b 1
)

pg_isready -h localhost -p 5432 >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: PostgreSQL is not running on localhost:5432
    echo Please start PostgreSQL service first
    echo Alternative: Use Docker - docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=@Salman2021 postgres:15
    pause
    exit /b 1
)

echo ✓ Go and PostgreSQL are available

REM Step 2: Setup Environment
echo [2/6] Setting up environment...
if not exist .env (
    copy .env.example .env
    echo ✓ Created .env file from template
) else (
    echo ✓ .env file already exists
)

REM Step 3: Download Dependencies
echo [3/6] Downloading Go dependencies...
go mod tidy
if %errorlevel% neq 0 (
    echo ERROR: Failed to download dependencies
    pause
    exit /b 1
)
echo ✓ Dependencies downloaded

REM Step 4: Setup Database
echo [4/6] Setting up database...
call scripts\setup-database.bat
if %errorlevel% neq 0 (
    echo ERROR: Database setup failed
    pause
    exit /b 1
)
echo ✓ Database setup completed

REM Step 5: Build Application
echo [5/6] Building application...
go build -o main.exe cmd\server\main.go
if %errorlevel% neq 0 (
    echo ERROR: Build failed
    pause
    exit /b 1
)
echo ✓ Application built successfully

REM Step 6: Start Application
echo [6/6] Starting application...
echo.
echo ========================================
echo 🚀 LASANI Geopolitical Service Starting
echo ========================================
echo.
echo Access Points:
echo • Application: http://localhost:8081
echo • Health Check: http://localhost:8081/health
echo • Countries API: http://localhost:8081/api/v1/countries
echo • Business Central UI: http://localhost:8081/web-ui/business-central.html
echo.
echo Press Ctrl+C to stop the service
echo ========================================
echo.

main.exe