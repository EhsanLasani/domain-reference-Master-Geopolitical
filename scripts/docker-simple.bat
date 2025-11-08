@echo off
echo Starting essential services only...

REM Check if Docker is running
docker --version >nul 2>&1
if errorlevel 1 (
    echo ERROR: Docker is not running or not installed
    echo Please start Docker Desktop and try again
    pause
    exit /b 1
)

echo Starting PostgreSQL, Redis, Jaeger, and Grafana...
docker-compose -f docker-compose.simple.yml up -d

echo.
echo Waiting for services to start...
timeout /t 10 /nobreak >nul

echo.
echo Services should be available at:
echo - PostgreSQL: localhost:5432
echo - Redis: localhost:6379
echo - Jaeger UI: http://localhost:16686
echo - Grafana: http://localhost:3000 (admin/admin123)
echo.
echo Check status: docker-compose -f docker-compose.simple.yml ps
echo View logs: docker-compose -f docker-compose.simple.yml logs -f