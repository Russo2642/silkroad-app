package pg

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"silkroad/m/internal/domain/forms"
)

type HelpWithTourFormPostgres struct {
	db *sqlx.DB
}

func NewHelpWithTourForm(db *sqlx.DB) *HelpWithTourFormPostgres {
	return &HelpWithTourFormPostgres{db: db}
}

func (r *HelpWithTourFormPostgres) Create(helpWithTourForm forms.HelpWithTourForm) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createContactFormQuery := fmt.Sprintf("INSERT INTO %s (name, phone, place, when_date)"+
		" VALUES ($1, $2, $3, $4) RETURNING id",
		helpWithTourFormTable)
	row := tx.QueryRow(createContactFormQuery, helpWithTourForm.Name, helpWithTourForm.Phone, helpWithTourForm.Place, helpWithTourForm.WhenDate)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}
