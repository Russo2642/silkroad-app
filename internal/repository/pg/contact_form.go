package pg

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"silkroad/m/internal/domain/forms"
)

type ContactFormPostgres struct {
	db *sqlx.DB
}

func NewContactForm(db *sqlx.DB) *ContactFormPostgres {
	return &ContactFormPostgres{db: db}
}

func (r *ContactFormPostgres) Create(contactForm forms.ContactForm) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createContactFormQuery := fmt.Sprintf("INSERT INTO %s (name, phone, email, description, tour_id)"+
		" VALUES ($1, $2, $3, $4, $5) RETURNING id",
		contactFormTable)
	row := tx.QueryRow(createContactFormQuery, contactForm.Name, contactForm.Phone, contactForm.Email,
		contactForm.Description, contactForm.TourID)
	if err := row.Scan(&id); err != nil {
		err := tx.Rollback()
		if err != nil {
			return 0, err
		}
		return 0, err
	}

	return id, tx.Commit()
}
