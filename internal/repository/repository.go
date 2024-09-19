package repository

import (
	"github.com/jmoiron/sqlx"
	"silkroad/m/internal/domain/forms"
	"silkroad/m/internal/domain/tour"
	"silkroad/m/internal/repository/pg"
)

type ContactForm interface {
	Create(contactForm forms.ContactForm) (int, error)
}

type HelpWithTourForm interface {
	Create(helpWithTourForm forms.HelpWithTourForm) (int, error)
}

type Tour interface {
	Create(tour tour.Tour) (int, error)
	GetAll(priceRange, tourPlace, tourDate, searchTitle string, quantity, duration int, limit, offset int) ([]tour.Tour, error)
	GetById(tourId int) (tour.Tour, error)
}

type TourEditor interface {
	Create(tourEditor tour.TourEditor) (int, error)
}

type Repository struct {
	ContactForm
	HelpWithTourForm
	Tour
	TourEditor
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		ContactForm:      pg.NewContactForm(db),
		HelpWithTourForm: pg.NewHelpWithTourForm(db),
		Tour:             pg.NewTour(db),
		TourEditor:       pg.NewTourEditor(db),
	}
}
