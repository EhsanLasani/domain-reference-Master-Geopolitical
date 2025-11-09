@echo off
echo Testing API endpoints...

echo [1] Health check:
curl -s http://localhost:8081/health
echo.

echo [2] Countries API:
curl -s http://localhost:8081/api/v1/countries
echo.

echo [3] Countries API with verbose:
curl -v http://localhost:8081/api/v1/countries
echo.

pause