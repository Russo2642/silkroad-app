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

func (s *TourService) GetAll(priceRange, tourPlace, tourDate, searchTitle string, quantity []int, duration, limit, offset int) ([]tour.Tour, int, error) {
	return s.repo.GetAll(priceRange, tourPlace, tourDate, searchTitle, quantity, duration, limit, offset)
}

func (s *TourService) GetById(tourId int) (tour.Tour, error) {
	return s.repo.GetById(tourId)
}

func (s *TourService) GetBySlug(tourSlug string) (tour.Tour, error) {
	return s.repo.GetBySlug(tourSlug)
}

func (s *TourService) GetMinMaxPrice() (int, int, error) {
	return s.repo.GetMinMaxPrice()
}
