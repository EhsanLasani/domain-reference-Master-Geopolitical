# Data Access Layer - Complete Implementation

## 🏗️ **Architecture Overview**

The data access layer implements comprehensive enterprise patterns with complete validation, error handling, events, and business catalogue.

```
📊 Database Schema (PostgreSQL)
    ↓
🔍 Validation Layer (Pre-database validation)
    ↓
🗄️ Repository Layer (CRUD + Business operations)
    ↓
❌ Error Mapping (Database → Domain errors)
    ↓
📡 Event System (Outbox pattern)
    ↓
📋 Business Catalogue (Available operations)
```

## 📁 **Directory Structure**

```
data-access-layer/
├── orm-odm-abstractions/
│   ├── country_model.go          # Domain models with validation tags
│   └── generated_models.go       # Auto-generated from schema
├── repositories-daos/
│   ├── country_repository_proper.go  # Manual repository
│   └── generated_repositories.go     # Auto-generated repositories
├── validation/
│   └── country_validator.go      # Pre-database validation
├── error-handling/
│   └── error_mapper.go           # Database → Domain error mapping
├── view-mappers/
│   └── view_mappers.go           # Go structs for database views
├── catalogue/
│   └── repository_catalogue.go   # Business layer operations catalogue
└── README.md
```

## ✅ **Implementation Checklist Compliance**

### **Database Schema Aspects Covered**

- ✅ **Validation**: Complete pre-database validation in `validation/`
- ✅ **Views**: Optimized query views in `views/database_views.sql`
- ✅ **Indexes**: Referenced in generated repositories for performance
- ✅ **Constraints**: Validated in code before database operations
- ✅ **Error Mapping**: PostgreSQL errors → Domain-specific error codes
- ✅ **LASANI Compliance**: 27-field audit system validation

### **Business Layer Catalogue**

The `catalogue/repository_catalogue.go` provides:

```go
// Available operations for business layer
type CountryOperations struct {
    // CRUD Operations
    Create, GetByID, GetByCode, Update, Delete
    
    // Query Operations  
    ListActive, Search
    
    // Business Operations
    ExistsByCode, GetActiveCount
    
    // Validation Operations
    ValidateCode, ValidateName
}
```

### **Event System Organization**

Following enterprise guidelines with outbox pattern:

```go
// Domain Events
CountryCreatedEvent    = "geo.country.created.v1"
CountryUpdatedEvent    = "geo.country.updated.v1"
CountryDeactivatedEvent = "geo.country.deactivated.v1"

// Outbox Pattern
type OutboxEvent struct {
    EventType, AggregateID, EventData
    Status: pending/processed/failed
    RetryCount, ErrorMsg
}
```

### **Error Handling in All Layers**

**Domain Error Codes:**
- `GEO-1xxx`: Country errors (1001=NotFound, 1002=Duplicate, 1003=Invalid, 1004=VersionConflict)
- `GEO-2xxx`: Region errors
- `GEO-3xxx`: Language errors  
- `GEO-9xxx`: System errors

**Error Flow:**
```
Database Error → ErrorMapper → DomainError → Business Layer → Presentation Layer
```

## 🔄 **Auto-Generation System**

### **Schema-to-Code Generator**
- **File**: `tools/schema-to-code/main.go`
- **Generates**: Models, repositories, validation rules
- **Triggers**: Schema changes, manual execution

### **Schema Watcher**
- **File**: `tools/schema-watcher.go`  
- **Monitors**: Real-time schema changes via MD5 checksum
- **Auto-regenerates**: Code when schema changes detected

### **Usage**
```bash
# Manual generation
go run tools/schema-to-code/main.go domain_reference_master_geopolitical

# Auto-watch mode
go run tools/schema-watcher.go

# Batch update
scripts\update-data-access-layer.bat
```

## 📊 **Database Views Available**

- `v_countries_active`: Active countries only (most used)
- `v_countries_with_regions`: Countries with region joins
- `v_audit_trail`: Complete audit history for compliance

## 🎯 **Business Layer Integration**

Business layer receives:

1. **Operations Catalogue**: All available functions with metadata
2. **Validation Results**: Pre-validated data before database
3. **Domain Events**: Structured events for business logic
4. **Error Codes**: Standardized error handling
5. **Performance Metadata**: Cache TTL, operation costs

## 🚀 **Performance Features**

- **Indexed Queries**: All operations use database indexes
- **View Optimization**: Pre-joined data for common queries  
- **Validation Caching**: In-memory validation rules
- **Connection Pooling**: Efficient database connection management
- **Error Caching**: Prevent repeated validation failures

This implementation provides complete enterprise-grade data access with validation, error handling, events, and business integration following all implementation checklist requirements.