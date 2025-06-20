package pg

import (
	"database/sql"
	"fmt"

	"silkroad/m/internal/domain/forms"

	"github.com/jmoiron/sqlx"
)

type HelpWithTourFormRepository struct {
	db *sqlx.DB
}

func NewHelpWithTourForm(db *sqlx.DB) *HelpWithTourFormRepository {
	return &HelpWithTourFormRepository{db: db}
}

func (r *HelpWithTourFormRepository) Create(helpWithTourForm forms.HelpWithTourForm) (int, error) {
	var id int
	query := `
		INSERT INTO help_with_tour_form (name, phone, country, when_date, created_at) 
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) 
		RETURNING id`

	err := r.db.QueryRow(query,
		helpWithTourForm.Name, helpWithTourForm.Phone, helpWithTourForm.Country, helpWithTourForm.WhenDate,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create help with tour form: %w", err)
	}

	return id, nil
}

func (r *HelpWithTourFormRepository) GetByID(id int) (forms.HelpWithTourForm, error) {
	var hf forms.HelpWithTourForm
	query := `
		SELECT id, name, phone, country, when_date, created_at
		FROM help_with_tour_form 
		WHERE id = $1`

	err := r.db.Get(&hf, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return hf, fmt.Errorf("help with tour form not found")
		}
		return hf, fmt.Errorf("failed to get help with tour form: %w", err)
	}

	return hf, nil
}
