@echo off
echo Checking service status...

echo.
echo === Docker Services ===
docker-compose -f docker-compose.simple.yml ps

echo.
echo === Port Status ===
netstat -an | findstr ":3000 :16686 :5432 :6379"

echo.
echo === Service Health ===
echo Testing Jaeger...
curl -s http://localhost:16686 >nul 2>&1 && echo Jaeger: OK || echo Jaeger: NOT ACCESSIBLE

echo Testing Grafana...
curl -s http://localhost:3000 >nul 2>&1 && echo Grafana: OK || echo Grafana: NOT ACCESSIBLE

echo.
echo If services show as NOT ACCESSIBLE:
echo 1. Run: scripts\docker-simple.bat
echo 2. Wait 30 seconds for services to start
echo 3. Check Docker Desktop is running