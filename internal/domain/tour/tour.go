package tour

import (
	"time"
)

type TourType string

type DescriptionRoute struct {
	Default []string `json:"default" db:"default"`
	Next    []string `json:"next" db:"next"`
	Photos  []string `json:"photos" db:"photos"`
}

const (
	OneDayTour    TourType = "Однодневный тур"
	MultiDayTour  TourType = "Многодневный тур"
	CityTour      TourType = "Сити-тур"
	ExclusiveTour TourType = "Эксклюзивный тур"
	InfoTour      TourType = "Инфо-тур"
	AuthorsTour   TourType = "Авторский тур"
)

type Tour struct {
	Id                   int              `json:"id" db:"id"`
	TourType             TourType         `json:"tour_type" db:"tour_type" binding:"required"`
	Slug                 string           `json:"slug" db:"slug"`
	Title                string           `json:"title" db:"title" binding:"required"`
	TourPlace            string           `json:"tour_place" db:"tour_place" binding:"required"`
	Season               string           `json:"season" db:"season" binding:"required"`
	Quantity             int              `json:"quantity" db:"quantity" binding:"required"`
	Duration             int              `json:"duration" db:"duration" binding:"required"`
	PhysicalRating       int              `json:"physical_rating" db:"physical_rating" binding:"required,gte=1,lte=5"`
	DescriptionExcursion string           `json:"description_excursion" db:"description_excursion" binding:"required"`
	DescriptionRoute     DescriptionRoute `json:"description_route" db:"description_route" binding:"required"`
	Price                int              `json:"price" db:"price" binding:"required"`
	Currency             string           `json:"currency" db:"currency" binding:"required"`
	Activity             []string         `json:"activity" db:"activity" binding:"required"`
	Tariff               string           `json:"tariff" db:"tariff"`
	TourDate             time.Time        `json:"tour_date" db:"tour_date" binding:"required"`
	Photos               []string         `json:"photos" db:"photos"`
}

func IsValidTourType(t TourType) bool {
	switch t {
	case OneDayTour, MultiDayTour, CityTour, ExclusiveTour, InfoTour, AuthorsTour:
		return true
	}
	return false
}
