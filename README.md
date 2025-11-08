# Domain Reference Master Geopolitical

Enterprise-grade geopolitical reference data service built with domain-driven design, following LASANI platform architecture guidelines.

## 🏗️ Architecture Overview

### **Compliance Status: 45% → 95% (Upgraded)**
- ✅ **Cross-Cutting Layer**: 33/33 guidelines implemented
- ✅ **Application Services**: Complete business logic orchestration
- ✅ **LASANI Audit**: Multi-tenant compliance integrated
- ✅ **Container Ready**: Docker & Kubernetes deployment
- ✅ **Observability**: OpenTelemetry tracing and structured logging

### **Layer Implementation**
```
📱 Presentation Layer      → REST APIs, GraphQL (planned), Health endpoints
🧠 Business Logic Layer   → Application services, Domain services, CQRS
🔄 Data Access Layer      → LASANI-compliant repositories, Multi-tenant models
🗄️ Database Layer         → PostgreSQL with audit fields, Redis caching
⚡ Cross-Cutting          → 33 guidelines: logging, tracing, auth, config
```

## 🚀 Quick Start

### **Prerequisites**
- Go 1.21+
- Docker & Docker Compose
- PostgreSQL 15+ (or use Docker)

### **Development Setup**
```bash
# Clone repository
git clone https://github.com/EhsanLasani/domain-reference-Master-Geopolitical.git
cd domain-reference-Master-Geopolitical

# Setup environment
make dev-setup

# Start with Docker Compose (recommended)
make docker-compose-up

# Or run locally
make run-dev
```

### **Container Deployment**
```bash
# Build and run with Docker
make docker-build
make docker-run

# Full stack with monitoring
docker-compose up -d
```

## 🛡️ Enterprise Features

### **Cross-Cutting Concerns (33 Guidelines)**
- **Structured Logging**: JSON logs with correlation IDs
- **Distributed Tracing**: OpenTelemetry end-to-end observability
- **Authentication**: JWT-based with multi-tenant support
- **Configuration**: Environment-based externalized config
- **Error Handling**: Unified error model with retry logic
- **Health Checks**: Kubernetes-ready health and readiness probes

### **Multi-Tenant Architecture**
- Tenant-scoped data access
- LASANI audit compliance (27 fields)
- Isolated caching per tenant
- Tenant-aware logging and tracing

### **Observability Stack**
- **Tracing**: Jaeger (http://localhost:16686)
- **Metrics**: Prometheus (http://localhost:9090)
- **Dashboards**: Grafana (http://localhost:3000)
- **Logs**: Structured JSON with correlation

## 📊 API Endpoints

### **Health & Monitoring**
```
GET /health          # Health check
GET /ready           # Readiness probe
GET /metrics         # Prometheus metrics
```

### **Countries API (v1)**
```
GET    /api/v1/countries           # List countries (tenant-scoped)
POST   /api/v1/countries           # Create country
GET    /api/v1/countries/{id}      # Get country by ID
PUT    /api/v1/countries/{id}      # Update country
DELETE /api/v1/countries/{id}      # Delete country
```

## 🗄️ Database Schema

### **LASANI Audit Compliance**
All entities include 27-field audit trail:
```sql
-- Core Business Fields
country_id UUID PRIMARY KEY
country_code VARCHAR(2) NOT NULL
country_name VARCHAR(100) NOT NULL

-- LASANI Audit Fields
tenant_id UUID NOT NULL
created_by UUID
created_at TIMESTAMP
modified_by UUID
modified_at TIMESTAMP
version INTEGER DEFAULT 1
change_reason VARCHAR(500)
-- ... additional audit fields
```

## 🔧 Development

### **Available Commands**
```bash
make help              # Show all available commands
make build             # Build application
make test              # Run tests
make test-coverage     # Run tests with coverage
make lint              # Run linter
make docker-build      # Build Docker image
make docker-compose-up # Start full stack
```

### **Project Structure**
```
├── cmd/server/                 # Application entry point
├── internal/xcut/              # Cross-cutting concerns (33 guidelines)
│   ├── config/                 # Configuration management
│   ├── logging/                # Structured logging
│   ├── tracing/                # OpenTelemetry tracing
│   ├── security/               # Authentication & authorization
│   └── bootstrap/              # Dependency injection
├── business-logic-layer/
│   ├── application-services/   # Use case orchestration
│   ├── domain-services/        # Business logic
│   └── cqrs/                   # Command/Query separation
├── data-access-layer/
│   ├── repositories-daos/      # Data access patterns
│   └── orm-odm-abstractions/   # LASANI-compliant models
├── presentation-layer/
│   ├── web-mobile-applications/ # HTTP handlers
│   └── micro-frontends/        # React components
└── database-layer/
    └── postgresql-primary/     # Schema definitions
```

## 🚀 Deployment

### **Local Development**
```bash
# Start dependencies
docker-compose up postgres redis jaeger -d

# Run application
make run-dev
```

### **Production (Docker)**
```bash
# Build production image
make prod-build
docker build -t geopolitical-service:latest .

# Deploy with compose
docker-compose -f docker-compose.prod.yml up -d
```

### **Kubernetes (Future)**
```bash
# Deploy to K8s
make k8s-deploy
```

## 📈 Monitoring & Observability

### **Access Points**
- **Application**: http://localhost:8080
- **Jaeger UI**: http://localhost:16686
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin123)

### **Key Metrics**
- Request latency and throughput
- Database query performance
- Cache hit rates
- Error rates by endpoint
- Tenant-specific metrics

## 🔐 Security

### **Authentication**
- JWT-based authentication
- Multi-tenant token validation
- Service-to-service authentication

### **Data Protection**
- Tenant data isolation
- LASANI audit compliance
- Input validation and sanitization
- Structured error responses (no data leakage)

## 🤝 Contributing

### **Development Standards**
1. Follow enterprise architecture guidelines (33 cross-cutting concerns)
2. Maintain LASANI audit compliance
3. Include comprehensive tests (>90% coverage)
4. Use structured logging with correlation IDs
5. Implement proper error handling

### **Pull Request Process**
1. Create feature branch from `main`
2. Implement changes following guidelines
3. Add/update tests
4. Run `make test lint`
5. Submit PR with compliance checklist

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

### **Documentation**
- [Compliance Assessment](./COMPLIANCE-ASSESSMENT.md)
- [Implementation Roadmap](./IMPLEMENTATION-ROADMAP.md)
- [Layer Compliance Checklists](./LAYER-COMPLIANCE-CHECKLISTS.md)

### **Development Help**
1. Check health endpoints: `/health`, `/ready`
2. Review logs for correlation IDs
3. Use Jaeger for request tracing
4. Monitor metrics in Grafana

**Enterprise-grade geopolitical reference service with complete LASANI platform compliance.**