// Package database implements database connection management with pooling
package database

import (
	"fmt"

	"github.com/EhsanLasani/domain-reference-Master-Geopolitical/internal/xcut/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode)


	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxConnections)
	sqlDB.SetMaxIdleConns(cfg.MaxConnections / 2)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return &Database{DB: db}, nil
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (d *Database) AutoMigrate(models ...interface{}) error {
	// Handle tenant_id migration for all tables
	tables := []string{"countries", "regions", "languages", "timezones", "country_subdivisions", "locales"}
	
	for _, table := range tables {
		// Check if table exists
		var tableExists bool
		d.DB.Raw(`SELECT EXISTS (
			SELECT 1 FROM information_schema.tables 
			WHERE table_schema = 'domain_reference_master_geopolitical' 
			AND table_name = ?
		)`, table).Scan(&tableExists)
		
		if !tableExists {
			continue
		}
		
		// Check if tenant_id column exists
		var columnExists bool
		d.DB.Raw(`SELECT EXISTS (
			SELECT 1 FROM information_schema.columns 
			WHERE table_schema = 'domain_reference_master_geopolitical' 
			AND table_name = ? 
			AND column_name = 'tenant_id'
		)`, table).Scan(&columnExists)
		
		if !columnExists {
			// Add column as nullable first
			d.DB.Exec(fmt.Sprintf("ALTER TABLE domain_reference_master_geopolitical.%s ADD COLUMN tenant_id varchar(100)", table))
		}
		
		// Update all NULL/empty values
		d.DB.Exec(fmt.Sprintf("UPDATE domain_reference_master_geopolitical.%s SET tenant_id = 'default-tenant' WHERE tenant_id IS NULL OR tenant_id = ''", table))
		
		// Make column NOT NULL after updating values
		d.DB.Exec(fmt.Sprintf("ALTER TABLE domain_reference_master_geopolitical.%s ALTER COLUMN tenant_id SET NOT NULL", table))
	}
	
	return d.DB.AutoMigrate(models...)
}