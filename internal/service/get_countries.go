package service

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"silkroad/m/internal/config"
)

type CountryService struct {
	countries []string `yaml:"countries"`
}

func NewCountryService() *CountryService {
	return &CountryService{}
}

func (s *CountryService) GetAll() ([]string, error) {
	if err := config.InitConfig("configs", "countries"); err != nil {
		logrus.Fatalf("init config err: %s", err.Error())
	}

	countries := viper.GetStringSlice("countries")
	if len(countries) == 0 {
		logrus.Error("No countries found")
		return nil, nil
	}

	return countries, nil
}
