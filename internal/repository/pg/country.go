package pg

import (
	"database/sql"
	"fmt"
	"strings"

	"silkroad/m/internal/domain/country"

	"github.com/jmoiron/sqlx"
)

type CountryRepository struct {
	db *sqlx.DB
}

func NewCountryRepository(db *sqlx.DB) *CountryRepository {
	return &CountryRepository{db: db}
}

func (r *CountryRepository) Create(c country.Country) (int, error) {
	var id int
	query := `
		INSERT INTO countries (name, code, is_active) 
		VALUES ($1, $2, $3) 
		RETURNING id`

	err := r.db.QueryRow(query, c.Name, c.Code, c.IsActive).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create country: %w", err)
	}

	return id, nil
}

func (r *CountryRepository) GetByID(id int) (country.Country, error) {
	var c country.Country
	query := `
		SELECT id, name, code, is_active, created_at, updated_at 
		FROM countries 
		WHERE id = $1`

	err := r.db.Get(&c, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, fmt.Errorf("country not found")
		}
		return c, fmt.Errorf("failed to get country: %w", err)
	}

	return c, nil
}

func (r *CountryRepository) GetByCode(code string) (country.Country, error) {
	var c country.Country
	query := `
		SELECT id, name, code, is_active, created_at, updated_at 
		FROM countries 
		WHERE code = $1`

	err := r.db.Get(&c, query, code)
	if err != nil {
		if err == sql.ErrNoRows {
			return c, fmt.Errorf("country not found")
		}
		return c, fmt.Errorf("failed to get country: %w", err)
	}

	return c, nil
}

func (r *CountryRepository) GetAll(filter country.CountryFilter) ([]country.Country, error) {
	var countries []country.Country
	var conditions []string
	var args []interface{}
	argIndex := 1

	query := `
		SELECT id, name, code, is_active, created_at, updated_at 
		FROM countries`

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argIndex))
		args = append(args, *filter.IsActive)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY name"

	err := r.db.Select(&countries, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get countries: %w", err)
	}

	return countries, nil
}

func (r *CountryRepository) Update(c country.Country) error {
	query := `
		UPDATE countries 
		SET name = $1, code = $2, is_active = $3, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $4`

	result, err := r.db.Exec(query, c.Name, c.Code, c.IsActive, c.ID)
	if err != nil {
		return fmt.Errorf("failed to update country: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("country not found")
	}

	return nil
}

func (r *CountryRepository) Delete(id int) error {
	query := `DELETE FROM countries WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete country: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("country not found")
	}

	return nil
}

func (r *CountryRepository) GetActiveCountries() ([]country.Country, error) {
	filter := country.CountryFilter{
		IsActive: &[]bool{true}[0],
	}
	return r.GetAll(filter)
}
