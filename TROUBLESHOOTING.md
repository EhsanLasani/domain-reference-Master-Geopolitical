# Troubleshooting Guide

## Cannot Access Jaeger (localhost:16686) or Grafana (localhost:3000)

### **Quick Fix**
```cmd
# 1. Start Docker Desktop first
# 2. Run simplified services
scripts\docker-simple.bat

# 3. Check status
scripts\check-services.bat

# 4. Wait 30 seconds, then try again
```

### **Step-by-Step Diagnosis**

#### 1. Check Docker Desktop
- Open Docker Desktop application
- Ensure it's running (green icon in system tray)
- If not installed: Download from https://docker.com/products/docker-desktop

#### 2. Verify Docker is Working
```cmd
docker --version
docker ps
```

#### 3. Start Services
```cmd
# Use simplified compose file
scripts\docker-simple.bat

# Check if containers are running
docker-compose -f docker-compose.simple.yml ps
```

#### 4. Check Port Conflicts
```cmd
# Check if ports are in use
netstat -an | findstr ":3000 :16686"

# If ports are busy, stop conflicting services
```

#### 5. View Container Logs
```cmd
# Check Jaeger logs
docker logs geopolitical-jaeger

# Check Grafana logs  
docker logs geopolitical-grafana
```

### **Alternative: Run Without Docker**

#### Local Jaeger (if Docker fails)
```cmd
# Download Jaeger binary
# https://github.com/jaegertracing/jaeger/releases
jaeger-all-in-one.exe --collector.otlp.enabled=true
```

#### Local Development (No Monitoring)
```cmd
# Just run the Go application
scripts\run-dev.bat

# Access: http://localhost:8080/health
```

### **Common Issues**

#### "Docker not found"
- Install Docker Desktop
- Restart terminal after installation
- Add Docker to PATH

#### "Port already in use"
```cmd
# Find process using port
netstat -ano | findstr :3000
netstat -ano | findstr :16686

# Kill process (replace PID)
taskkill /PID <process_id> /F
```

#### "Container won't start"
```cmd
# Remove old containers
docker-compose -f docker-compose.simple.yml down -v

# Restart
scripts\docker-simple.bat
```

### **Minimal Setup (No Monitoring)**

If you just want to run the application:

```cmd
# 1. Setup environment
scripts\dev-setup.bat

# 2. Run application only
go run cmd\server\main.go

# 3. Test
curl http://localhost:8080/health
```

Access points:
- **Application**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **API**: http://localhost:8080/api/v1/countries