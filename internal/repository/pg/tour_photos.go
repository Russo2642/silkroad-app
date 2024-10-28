package pg

import (
	"fmt"
	"github.com/lib/pq"
	"mime/multipart"
	"silkroad/m/internal/aws"
	"time"
)

func (r *TourPostgres) AddPhotos(tourID int, files []*multipart.FileHeader) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var photoUrls []string

	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		key := fmt.Sprintf("tour_photos/%d_%s", time.Now().Unix(), file.Filename)

		fileUrl, err := aws.UploadPhotoToS3("gosilkroadbucket", key, f)
		if err != nil {
			_ = f.Close()
			return fmt.Errorf("failed to upload photo: %w", err)
		}

		photoUrls = append(photoUrls, fileUrl)

		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}
	}

	updatePhotosQuery := fmt.Sprintf("UPDATE %s SET photos = array_cat(photos, $1) WHERE id = $2", tourTable)
	_, err = tx.Exec(updatePhotosQuery, pq.Array(photoUrls), tourID)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}