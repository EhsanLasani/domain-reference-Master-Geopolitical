// ============================================================================
// FILE: view_mappers.go
// DOMAIN: Reference Master Geopolitical
// LAYER: Data Access Layer - View Mappers
// PURPOSE: Map database views to Go structs (no SQL here)
// VERSION: 1.0.0
// CREATED: 2025-11-07
// ============================================================================

package viewmappers

import (
	"database/sql"
)

// CountryActiveView maps v_countries_active database view
type CountryActiveView struct {
	CountryID    string `db:"country_id" json:"country_id"`
	CountryCode  string `db:"country_code" json:"country_code"`
	CountryName  string `db:"country_name" json:"country_name"`
	ISO3Code     string `db:"iso3_code" json:"iso3_code"`
	OfficialName string `db:"official_name" json:"official_name"`
}

// ViewRepository handles database view queries
type ViewRepository struct {
	db *sql.DB
}

func NewViewRepository(db *sql.DB) *ViewRepository {
	return &ViewRepository{db: db}
}

// GetActiveCountries queries v_countries_active view
func (r *ViewRepository) GetActiveCountries() ([]CountryActiveView, error) {
	query := "SELECT country_id, country_code, country_name, iso3_code, official_name FROM domain_reference_master_geopolitical.v_countries_active"
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []CountryActiveView
	for rows.Next() {
		var country CountryActiveView
		err := rows.Scan(&country.CountryID, &country.CountryCode, &country.CountryName, &country.ISO3Code, &country.OfficialName)
		if err != nil {
			continue
		}
		results = append(results, country)
	}
	
	return results, nil
}