package pg

import (
	"database/sql"
	"fmt"

	"silkroad/m/internal/domain/forms"

	"github.com/jmoiron/sqlx"
)

type ContactFormRepository struct {
	db *sqlx.DB
}

func NewContactForm(db *sqlx.DB) *ContactFormRepository {
	return &ContactFormRepository{db: db}
}

func (r *ContactFormRepository) Create(contactForm forms.ContactForm) (int, error) {
	var id int
	query := `
		INSERT INTO contact_form (name, phone, email, tour_id, created_at) 
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP) 
		RETURNING id`

	err := r.db.QueryRow(query,
		contactForm.Name, contactForm.Phone, contactForm.Email, contactForm.TourID,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create contact form: %w", err)
	}

	return id, nil
}

func (r *ContactFormRepository) GetByID(id int) (forms.ContactForm, error) {
	var cf forms.ContactForm
	query := `
		SELECT id, name, phone, email, tour_id, created_at
		FROM contact_form 
		WHERE id = $1`

	err := r.db.Get(&cf, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return cf, fmt.Errorf("contact form not found")
		}
		return cf, fmt.Errorf("failed to get contact form: %w", err)
	}

	return cf, nil
}
