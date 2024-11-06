package service

import (
	"fmt"
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

func (s *TourService) GetAll(tourPlace, tourDate, searchTitle string, quantity []int, priceMin, priceMax, duration, limit, offset int, popular bool) ([]tour.Tour, int, int, int, int, []string, error) {
	return s.repo.GetAll(tourPlace, tourDate, searchTitle, quantity, priceMin, priceMax, duration, limit, offset, popular)
}

func (s *TourService) GetById(tourId int) (tour.Tour, error) {
	return s.repo.GetTourByField("id", fmt.Sprintf("%d", tourId))
}

func (s *TourService) GetBySlug(tourSlug string) (tour.Tour, error) {
	return s.repo.GetTourByField("slug", tourSlug)
}

func (s *TourService) GetMinMaxPrice() (int, int, error) {
	return s.repo.GetMinMaxPrice()
}
