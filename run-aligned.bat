@echo off
echo ========================================
echo  ALIGNED SCHEMA APPLICATION STARTUP
echo ========================================
echo.
echo Database: geopolitical
echo Schema: domain_reference_master_geopolitical  
echo Port: 8082
echo Alignment: ✅ Complete
echo.
echo Starting aligned application...
echo.

go run cmd/server/main_aligned.go

pause