package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"github.com/google/uuid"
)

func main() {
	// Database connection
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=@Salman2021 dbname=referencemaster sslmode=disable")
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	defer db.Close()

	// Set schema
	_, err = db.Exec("SET search_path TO domain_reference_master_geopolitical, public")
	if err != nil {
		log.Fatal("Failed to set schema:", err)
	}

	fmt.Println("Seeding sample data...")

	// Seed Countries
	countries := []struct {
		code, name, iso3, official, capital, continent, phone string
		numeric                                               int
	}{
		{"US", "United States", "USA", "United States of America", "Washington D.C.", "NA", "+1", 840},
		{"GB", "United Kingdom", "GBR", "United Kingdom of Great Britain and Northern Ireland", "London", "EU", "+44", 826},
		{"DE", "Germany", "DEU", "Federal Republic of Germany", "Berlin", "EU", "+49", 276},
		{"FR", "France", "FRA", "French Republic", "Paris", "EU", "+33", 250},
		{"JP", "Japan", "JPN", "Japan", "Tokyo", "AS", "+81", 392},
		{"CN", "China", "CHN", "People's Republic of China", "Beijing", "AS", "+86", 156},
		{"IN", "India", "IND", "Republic of India", "New Delhi", "AS", "+91", 356},
		{"BR", "Brazil", "BRA", "Federative Republic of Brazil", "Brasília", "SA", "+55", 76},
		{"AU", "Australia", "AUS", "Commonwealth of Australia", "Canberra", "OC", "+61", 36},
		{"CA", "Canada", "CAN", "Canada", "Ottawa", "NA", "+1", 124},
	}

	for _, c := range countries {
		_, err = db.Exec(`
			INSERT INTO countries (country_id, country_code, country_name, iso3_code, numeric_code, 
				official_name, capital_city, continent_code, phone_prefix, is_active, is_deleted, 
				tenant_id, created_at, created_by, updated_at, updated_by, version)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, true, false, 'default-tenant', $10, 'system', $10, 'system', 1)
			ON CONFLICT (country_code, tenant_id) DO NOTHING`,
			uuid.New(), c.code, c.name, c.iso3, c.numeric, c.official, c.capital, c.continent, c.phone, time.Now())
	}
	fmt.Println("✓ Countries seeded")

	// Seed Regions
	regions := []struct {
		code, name, rtype string
		parent            *string
	}{
		{"NA", "North America", "CONTINENT", nil},
		{"EU", "Europe", "CONTINENT", nil},
		{"AS", "Asia", "CONTINENT", nil},
		{"SA", "South America", "CONTINENT", nil},
		{"AF", "Africa", "CONTINENT", nil},
		{"OC", "Oceania", "CONTINENT", nil},
	}

	for _, r := range regions {
		_, err = db.Exec(`
			INSERT INTO regions (region_id, region_code, region_name, region_type, parent_region_id,
				is_active, is_deleted, tenant_id, created_at, created_by, updated_at, updated_by, version)
			VALUES ($1, $2, $3, $4, $5, true, false, 'default-tenant', $6, 'system', $6, 'system', 1)
			ON CONFLICT (region_code, tenant_id) DO NOTHING`,
			uuid.New(), r.code, r.name, r.rtype, r.parent, time.Now())
	}
	fmt.Println("✓ Regions seeded")

	// Seed Languages
	languages := []struct {
		code, name, iso3, native, direction string
	}{
		{"en", "English", "eng", "English", "LTR"},
		{"es", "Spanish", "spa", "Español", "LTR"},
		{"fr", "French", "fra", "Français", "LTR"},
		{"de", "German", "deu", "Deutsch", "LTR"},
		{"zh", "Chinese", "zho", "中文", "LTR"},
		{"ja", "Japanese", "jpn", "日本語", "LTR"},
		{"ar", "Arabic", "ara", "العربية", "RTL"},
		{"hi", "Hindi", "hin", "हिन्दी", "LTR"},
		{"pt", "Portuguese", "por", "Português", "LTR"},
		{"ru", "Russian", "rus", "Русский", "LTR"},
	}

	for _, l := range languages {
		_, err = db.Exec(`
			INSERT INTO languages (language_id, language_code, language_name, iso3_code, native_name, direction,
				is_active, is_deleted, tenant_id, created_at, created_by, updated_at, updated_by, version)
			VALUES ($1, $2, $3, $4, $5, $6, true, false, 'default-tenant', $7, 'system', $7, 'system', 1)
			ON CONFLICT (language_code, tenant_id) DO NOTHING`,
			uuid.New(), l.code, l.name, l.iso3, l.native, l.direction, time.Now())
	}
	fmt.Println("✓ Languages seeded")

	// Seed Timezones
	timezones := []struct {
		code, name           string
		hours, minutes       int
		dst                  bool
		dstHours             *int
	}{
		{"UTC", "Coordinated Universal Time", 0, 0, false, nil},
		{"EST", "Eastern Standard Time", -5, 0, true, intPtr(1)},
		{"PST", "Pacific Standard Time", -8, 0, true, intPtr(1)},
		{"GMT", "Greenwich Mean Time", 0, 0, true, intPtr(1)},
		{"CET", "Central European Time", 1, 0, true, intPtr(1)},
		{"JST", "Japan Standard Time", 9, 0, false, nil},
		{"CST", "China Standard Time", 8, 0, false, nil},
		{"IST", "India Standard Time", 5, 30, false, nil},
		{"AEST", "Australian Eastern Standard Time", 10, 0, true, intPtr(1)},
		{"BRT", "Brasília Time", -3, 0, false, nil},
	}

	for _, tz := range timezones {
		_, err = db.Exec(`
			INSERT INTO timezones (timezone_id, timezone_code, timezone_name, utc_offset_hours, utc_offset_minutes,
				supports_dst, dst_offset_hours, is_active, is_deleted, tenant_id, created_at, created_by, updated_at, updated_by, version)
			VALUES ($1, $2, $3, $4, $5, $6, $7, true, false, 'default-tenant', $8, 'system', $8, 'system', 1)
			ON CONFLICT (timezone_code, tenant_id) DO NOTHING`,
			uuid.New(), tz.code, tz.name, tz.hours, tz.minutes, tz.dst, tz.dstHours, time.Now())
	}
	fmt.Println("✓ Timezones seeded")

	fmt.Println("Sample data seeding completed successfully!")
}

func intPtr(i int) *int {
	return &i
}