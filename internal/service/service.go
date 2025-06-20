package service

import (
	"mime/multipart"
	"silkroad/m/internal/domain/country"
	"silkroad/m/internal/domain/forms"
	"silkroad/m/internal/domain/tour"
	"silkroad/m/internal/repository"
)

type ContactForm interface {
	Create(contactForm forms.ContactForm) (int, error)
	GetByID(id int) (forms.ContactForm, error)
}

type HelpWithTourForm interface {
	Create(helpWithTourForm forms.HelpWithTourForm) (int, error)
	GetByID(id int) (forms.HelpWithTourForm, error)
}

type Country interface {
	Create(country country.Country) (int, error)
	GetByID(id int) (country.Country, error)
	GetByCode(code string) (country.Country, error)
	GetAll(filter country.CountryFilter) ([]country.Country, error)
	Update(country country.Country) error
	Delete(id int) error
	GetActiveCountries() ([]country.Country, error)
}

type Tour interface {
	Create(tour tour.Tour) (int, error)
	GetByID(id int) (tour.Tour, error)
	GetBySlug(slug string) (tour.Tour, error)
	GetAll(filter tour.TourFilter) ([]tour.Tour, int, error)
	GetSummaries(filter tour.TourFilter) ([]tour.TourSummary, int, error)
	Update(tour tour.Tour) error
	Delete(id int) error
	GetMinMaxPrice() (int, int, error)
	GetFilterValues() (map[string][]string, error)
}

type TourPhoto interface {
	Create(photo tour.TourPhotoInput, photoUrl string) (int, error)
	GetByID(id int) (tour.TourPhoto, error)
	GetByTourID(tourID int) (*tour.TourPhotosGrouped, error)
	GetByFilter(filter tour.TourPhotoFilter) ([]tour.TourPhoto, int, error)
	Update(id int, photo tour.TourPhotoInput) error
	Delete(id int) error
	DeleteByTourID(tourID int) error
	UploadPhotos(tourID int, files []*multipart.FileHeader, photoType tour.TourPhotoType, metadata tour.TourPhotoInput) error
}

type Service struct {
	ContactForm
	HelpWithTourForm
	Country
	Tour
	TourPhoto
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		ContactForm:      NewContactFormService(repos.ContactForm),
		HelpWithTourForm: NewHelpWithTourFormService(repos.HelpWithTourForm),
		Country:          NewCountryService(repos.Country),
		Tour:             NewTourService(repos.Tour),
		TourPhoto:        NewTourPhotoService(repos.TourPhoto),
	}
}
