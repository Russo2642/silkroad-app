package service

import (
	"silkroad/m/internal/domain/tour"
	"silkroad/m/internal/repository"
)

type TourEditorSerive struct {
	repo repository.TourEditor
}

func NewTourEditorService(repo repository.TourEditor) *TourEditorSerive {
	return &TourEditorSerive{repo: repo}
}

func (s *TourEditorSerive) Create(tourEditor tour.TourEditor) (int, error) {
	return s.repo.Create(tourEditor)
}
