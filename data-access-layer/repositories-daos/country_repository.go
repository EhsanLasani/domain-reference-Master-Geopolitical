// DATA ACCESS LAYER - Repository
package repositories

import (
	"database/sql"
	"reference-master-geopolitical/models"
)

type CountryRepository struct {
	db *sql.DB
}

func NewCountryRepository(db *sql.DB) *CountryRepository {
	return &CountryRepository{db: db}
}

func (r *CountryRepository) GetAllActiveCountries() ([]models.Country, error) {
	query := `SELECT country_code, country_name, COALESCE(iso3_code, ''), COALESCE(official_name, ''), is_active 
			  FROM domain_reference_master_geopolitical.countries 
			  WHERE is_active = true`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var country models.Country
		err := rows.Scan(&country.CountryCode, &country.CountryName, &country.ISO3Code, &country.OfficialName, &country.IsActive)
		if err != nil {
			continue
		}
		countries = append(countries, country)
	}
	return countries, nil
}