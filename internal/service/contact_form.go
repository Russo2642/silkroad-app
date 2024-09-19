package service

import (
	"silkroad/m/internal/domain/forms"
	"silkroad/m/internal/repository"
)

type ContactFormService struct {
	repo repository.ContactForm
}

func NewContactFormService(repo repository.ContactForm) *ContactFormService {
	return &ContactFormService{repo: repo}
}

func (s *ContactFormService) Create(contactForm forms.ContactForm) (int, error) {
	return s.repo.Create(contactForm)
}
