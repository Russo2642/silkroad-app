package service

import (
	"silkroad/m/internal/domain/country"
	"silkroad/m/internal/repository"
)

type CountryService struct {
	repo repository.Country
}

func NewCountryService(repo repository.Country) *CountryService {
	return &CountryService{repo: repo}
}

func (s *CountryService) Create(c country.Country) (int, error) {
	return s.repo.Create(c)
}

func (s *CountryService) GetByID(id int) (country.Country, error) {
	return s.repo.GetByID(id)
}

func (s *CountryService) GetByCode(code string) (country.Country, error) {
	return s.repo.GetByCode(code)
}

func (s *CountryService) GetAll(filter country.CountryFilter) ([]country.Country, error) {
	return s.repo.GetAll(filter)
}

func (s *CountryService) Update(c country.Country) error {
	return s.repo.Update(c)
}

func (s *CountryService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *CountryService) GetActiveCountries() ([]country.Country, error) {
	return s.repo.GetActiveCountries()
}
