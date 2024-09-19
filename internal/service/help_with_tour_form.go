package service

import (
	"silkroad/m/internal/domain/forms"
	"silkroad/m/internal/repository"
)

type HelpWithTourFormService struct {
	repo repository.HelpWithTourForm
}

func NewHelpWithTourFormService(repo repository.HelpWithTourForm) *HelpWithTourFormService {
	return &HelpWithTourFormService{repo: repo}
}

func (s *HelpWithTourFormService) Create(helpWithTourForm forms.HelpWithTourForm) (int, error) {
	return s.repo.Create(helpWithTourForm)
}
