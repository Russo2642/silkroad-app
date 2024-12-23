package repository

import (
	"github.com/jmoiron/sqlx"
	"mime/multipart"
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
	GetAll(tourPlace, tourDate, searchTitle string, quantity []int, priceMin, priceMax, duration, limit, offset int, popular bool) ([]tour.Tour, int, int, int, int, []string, error)
	GetTourByField(field, value string) (tour.Tour, error)
	GetMinMaxPrice() (int, int, error)
	AddPhotos(tourID int, files []*multipart.FileHeader, photoType string) error
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
