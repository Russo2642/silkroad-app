package service

import "mime/multipart"

func (s *TourService) AddPhotos(tourID int, files []*multipart.FileHeader, photoType string) error {
	return s.repo.AddPhotos(tourID, files, photoType)
}
