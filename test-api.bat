@echo off
echo Testing LASANI Geopolitical API...

REM Test 1: Health Check
echo [1/4] Testing health endpoint...
curl -s http://localhost:8081/health
echo.

REM Test 2: Get Countries
echo [2/4] Testing countries endpoint...
curl -s http://localhost:8081/api/v1/countries
echo.

REM Test 3: Create Country
echo [3/4] Testing create country...
curl -s -X POST http://localhost:8081/api/v1/countries ^
  -H "Content-Type: application/json" ^
  -d "{\"country_code\":\"TS\",\"country_name\":\"Test Country\"}"
echo.

REM Test 4: Get Specific Country
echo [4/4] Testing get country by code...
curl -s http://localhost:8081/api/v1/countries/TS
echo.

echo Testing completed!