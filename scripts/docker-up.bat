@echo off
echo Starting all services with Docker Compose...
docker-compose up -d

echo.
echo Services starting...
echo - Application: http://localhost:8080
echo - Jaeger UI: http://localhost:16686
echo - Prometheus: http://localhost:9090
echo - Grafana: http://localhost:3000
echo.
echo Run 'docker-compose logs -f' to see logs