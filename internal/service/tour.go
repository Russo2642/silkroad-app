package service

import (
	"time"

	"silkroad/m/internal/domain/tour"
	"silkroad/m/internal/repository"
)

type TourService struct {
	repo repository.Tour
}

func NewTourService(repo repository.Tour) *TourService {
	return &TourService{repo: repo}
}

func (s *TourService) Create(t tour.Tour) (int, error) {
	if t.Status == "" {
		t.Status = tour.StatusActive
	}

	if t.MinParticipants == 0 {
		t.MinParticipants = 1
	}

	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now

	return s.repo.Create(t)
}

func (s *TourService) GetByID(id int) (tour.Tour, error) {
	return s.repo.GetByID(id)
}

func (s *TourService) GetBySlug(slug string) (tour.Tour, error) {
	return s.repo.GetBySlug(slug)
}

func (s *TourService) GetAll(filter tour.TourFilter) ([]tour.Tour, int, error) {
	if filter.Limit == 0 {
		filter.Limit = 20
	}

	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}

	return s.repo.GetAll(filter)
}

func (s *TourService) GetSummaries(filter tour.TourFilter) ([]tour.TourSummary, int, error) {
	if filter.Limit == 0 {
		filter.Limit = 20
	}

	if filter.SortBy == "" {
		filter.SortBy = "created_at"
	}

	return s.repo.GetSummaries(filter)
}

func (s *TourService) Update(t tour.Tour) error {
	t.UpdatedAt = time.Now()
	return s.repo.Update(t)
}

func (s *TourService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *TourService) GetMinMaxPrice() (int, int, error) {
	return s.repo.GetMinMaxPrice()
}

func (s *TourService) GetFilterValues() (map[string][]string, error) {
	return s.repo.GetFilterValues()
}
