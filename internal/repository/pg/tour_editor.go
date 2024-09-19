package pg

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"silkroad/m/internal/domain/tour"
)

type TourEditorPostgres struct {
	db *sqlx.DB
}

func NewTourEditor(db *sqlx.DB) *TourEditorPostgres {
	return &TourEditorPostgres{db: db}
}

func (r *TourEditorPostgres) Create(tourEditor tour.TourEditor) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, nil
	}

	var id int
	createTourEditorQuery := fmt.Sprintf("INSERT INTO %s (name, phone, email, tour_date, activity, location)"+
		" VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", tourEditorTable)

	row := tx.QueryRow(createTourEditorQuery, tourEditor.Name, tourEditor.Phone, tourEditor.Email, tourEditor.TourDate,
		pq.Array(tourEditor.Activity), pq.Array(tourEditor.Location))
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}
