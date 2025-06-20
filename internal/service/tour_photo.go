package service

import (
	"fmt"
	"mime/multipart"
	"silkroad/m/internal/domain/tour"
	"silkroad/m/internal/repository"
)

type TourPhotoService struct {
	repo repository.TourPhoto
}

func NewTourPhotoService(repo repository.TourPhoto) *TourPhotoService {
	return &TourPhotoService{repo: repo}
}

func (s *TourPhotoService) Create(photo tour.TourPhotoInput, photoUrl string) (int, error) {
	return s.repo.Create(photo, photoUrl)
}

func (s *TourPhotoService) GetByID(id int) (tour.TourPhoto, error) {
	return s.repo.GetByID(id)
}

func (s *TourPhotoService) GetByTourID(tourID int) (*tour.TourPhotosGrouped, error) {
	return s.repo.GetByTourID(tourID)
}

func (s *TourPhotoService) GetByFilter(filter tour.TourPhotoFilter) ([]tour.TourPhoto, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	return s.repo.GetByFilter(filter)
}

func (s *TourPhotoService) Update(id int, photo tour.TourPhotoInput) error {
	return s.repo.Update(id, photo)
}

func (s *TourPhotoService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *TourPhotoService) DeleteByTourID(tourID int) error {
	return s.repo.DeleteByTourID(tourID)
}

func (s *TourPhotoService) UploadPhotos(tourID int, files []*multipart.FileHeader, photoType tour.TourPhotoType, metadata tour.TourPhotoInput) error {
	if len(files) == 0 {
		return fmt.Errorf("no files provided")
	}

	if !tour.IsValidPhotoType(photoType) {
		return fmt.Errorf("invalid photo type: %s", photoType)
	}

	maxFiles := 10
	if len(files) > maxFiles {
		return fmt.Errorf("too many files, maximum %d allowed", maxFiles)
	}

	maxFileSize := int64(10 << 20)
	for _, file := range files {
		if file.Size > maxFileSize {
			return fmt.Errorf("file %s is too large, maximum size is 10MB", file.Filename)
		}
	}

	metadata.TourID = tourID
	metadata.PhotoType = photoType

	return s.repo.UploadPhotos(tourID, files, photoType, metadata)
}
