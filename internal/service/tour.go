package service

import (
	"silkroad/m/internal/domain/tour"
	"silkroad/m/internal/repository"
)

type TourService struct {
	repo repository.Tour
}

func NewTourService(repo repository.Tour) *TourService {
	return &TourService{repo: repo}
}

func (s *TourService) Create(tour tour.Tour) (int, error) {
	return s.repo.Create(tour)
}

func (s *TourService) GetAll(priceRange, tourPlace, tourDate, searchTitle string, quantity, duration int, limit, offset int) ([]tour.Tour, error) {
	return s.repo.GetAll(priceRange, tourPlace, tourDate, searchTitle, quantity, duration, limit, offset)
}

func (s *TourService) GetById(tourId int) (tour.Tour, error) {
	return s.repo.GetById(tourId)
}
