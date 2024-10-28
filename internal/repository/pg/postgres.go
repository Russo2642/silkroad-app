package pg

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

const (
	contactFormTable      = "contact_form"
	helpWithTourFormTable = "help_with_tour_form"
	tourTable             = "tours"
	tourEditorTable       = "tour_editor"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s "+
		"sslmode=%s", cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
