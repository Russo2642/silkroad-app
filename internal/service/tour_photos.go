package service

import "mime/multipart"

func (s *TourService) AddPhotos(tourID int, files []*multipart.FileHeader, updateField string) error {
	return s.repo.AddPhotos(tourID, files, updateField)
}
