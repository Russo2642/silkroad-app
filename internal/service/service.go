package service

import (
	"mime/multipart"
	"silkroad/m/internal/domain/forms"
	"silkroad/m/internal/domain/tour"
	"silkroad/m/internal/repository"
)

type ContactForm interface {
	Create(contactForm forms.ContactForm) (int, error)
}

type HelpWithTourForm interface {
	Create(helpWithTourForm forms.HelpWithTourForm) (int, error)
}

type GetCountries interface {
	GetAll() ([]string, error)
}

type Tour interface {
	Create(tour tour.Tour) (int, error)
	GetAll(tourPlace, tourDate, searchTitle string, quantity []int, priceMin, priceMax, duration, limit, offset int) ([]tour.Tour, int, int, int, int, []string, error)
	GetById(tourId int) (tour.Tour, error)
	GetBySlug(tourSlug string) (tour.Tour, error)
	GetMinMaxPrice() (int, int, error)
	AddPhotos(tourID int, files []*multipart.FileHeader, updateField string) error
}

type TourEditor interface {
	Create(tourEditor tour.TourEditor) (int, error)
}

type Service struct {
	ContactForm
	HelpWithTourForm
	GetCountries
	Tour
	TourEditor
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		ContactForm:      NewContactFormService(repos.ContactForm),
		HelpWithTourForm: NewHelpWithTourFormService(repos.HelpWithTourForm),
		GetCountries:     NewCountryService(),
		Tour:             NewTourService(repos.Tour),
		TourEditor:       NewTourEditorService(repos.TourEditor),
	}
}
