package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

func main() {
	// Connect to geopolitical database
	dsn := "host=localhost port=5432 user=postgres password=@Salman2021 dbname=geopolitical sslmode=disable"
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open connection:", err)
	}
	defer db.Close()
	
	// Create ref schema
	_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS ref")
	if err != nil {
		fmt.Printf("Schema creation error (might exist): %v\n", err)
	} else {
		fmt.Println("✅ Created ref schema!")
	}
	
	// Create countries table with LASANI audit fields
	createTable := `
	CREATE TABLE IF NOT EXISTS ref.countries (
		country_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		country_code VARCHAR(2) UNIQUE NOT NULL,
		country_name VARCHAR(100) NOT NULL,
		iso3_code VARCHAR(3),
		numeric_code SMALLINT,
		official_name VARCHAR(200),
		capital_city VARCHAR(100),
		continent_code VARCHAR(2),
		region_code VARCHAR(10),
		currency_code VARCHAR(3),
		phone_prefix VARCHAR(10),
		is_active BOOLEAN DEFAULT true NOT NULL,
		is_deleted BOOLEAN DEFAULT false NOT NULL,
		created_at TIMESTAMP DEFAULT now(),
		created_by UUID,
		created_ip INET,
		created_device JSONB,
		created_session UUID,
		created_location JSONB,
		updated_at TIMESTAMP DEFAULT now(),
		updated_by UUID,
		updated_ip INET,
		updated_device JSONB,
		updated_session UUID,
		updated_location JSONB,
		deleted_at TIMESTAMP,
		deleted_by UUID,
		deleted_ip INET,
		deleted_device JSONB,
		deleted_session UUID,
		deleted_location JSONB,
		source_system VARCHAR(50) DEFAULT 'reference_master',
		change_reason TEXT,
		version INTEGER DEFAULT 1 NOT NULL
	)`
	
	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal("Failed to create countries table:", err)
	}
	
	fmt.Println("✅ Created countries table with LASANI audit fields!")
	
	// Insert sample data
	insertData := `
	INSERT INTO ref.countries (country_code, country_name, iso3_code, official_name) 
	VALUES 
		('US', 'United States', 'USA', 'United States of America'),
		('CA', 'Canada', 'CAN', 'Canada'),
		('GB', 'United Kingdom', 'GBR', 'United Kingdom of Great Britain and Northern Ireland')
	ON CONFLICT (country_code) DO NOTHING`
	
	_, err = db.Exec(insertData)
	if err != nil {
		fmt.Printf("Sample data insertion error: %v\n", err)
	} else {
		fmt.Println("✅ Inserted sample countries!")
	}
}