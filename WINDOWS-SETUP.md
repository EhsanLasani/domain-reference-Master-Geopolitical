# Windows Setup Guide

## Quick Start (No Make Required)

### 1. Development Setup
```cmd
# Setup development environment
scripts\dev-setup.bat

# Start development server
scripts\run-dev.bat
```

### 2. Docker Setup
```cmd
# Start all services
scripts\docker-up.bat

# Stop services
docker-compose down

# View logs
docker-compose logs -f
```

### 3. Manual Commands

#### Build Application
```cmd
go build -o bin\geopolitical-service.exe cmd\server\main.go
```

#### Run Tests
```cmd
go test -v .\...
```

#### Docker Build
```cmd
docker build -t geopolitical-service:latest .
```

#### Docker Run
```cmd
docker run -p 8080:8080 --env-file .env geopolitical-service:latest
```

## GitHub Setup

### 1. Create Repository
1. Go to https://github.com/EhsanLasani
2. Click "New repository"
3. Name: `domain-reference-Master-Geopolitical`
4. Make it public
5. Don't initialize with README (we have one)

### 2. Push Code
```cmd
# Check current status
git status

# Push to GitHub (after creating repository)
git push -u origin master
```

## Service URLs

After running `scripts\docker-up.bat`:

- **Application**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Jaeger Tracing**: http://localhost:16686
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin123)

## Troubleshooting

### Git Push Issues
If you get "failed to push" error:
1. Make sure you created the GitHub repository first
2. Repository must be named exactly: `domain-reference-Master-Geopolitical`
3. Try: `git push -u origin master --force` (first time only)

### Docker Issues
If Docker commands fail:
1. Make sure Docker Desktop is running
2. Check: `docker --version`
3. Try: `docker-compose --version`

### Go Issues
If Go commands fail:
1. Check: `go version` (should be 1.21+)
2. Run: `go mod tidy`
3. Check GOPATH: `go env GOPATH`