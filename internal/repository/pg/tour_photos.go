package pg

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/lib/pq"
	"mime/multipart"
	"silkroad/m/internal/aws"
	"time"
)

func (r *TourPostgres) AddPhotos(tourID int, files []*multipart.FileHeader, updateField string) error {
	var exists bool
	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", tourTable)
	err := r.db.QueryRow(checkQuery, tourID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check tour existence: %w", err)
	}

	if !exists {
		return fmt.Errorf("tour with id %d does not exist", tourID)
	}

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

	var updateQuery string
	if updateField == "gallery" {
		updateQuery = fmt.Sprintf("UPDATE %s SET photos = array_cat(COALESCE(photos, ARRAY[]::text[]), $1) WHERE id = $2", tourTable)
		_, err = tx.Exec(updateQuery, pq.Array(photoUrls), tourID)
	} else if updateField == "route" {
		descriptionRouteJSON, err := json.Marshal(photoUrls)
		if err != nil {
			return fmt.Errorf("failed to marshal photo URLs: %w", err)
		}
		updateQuery = fmt.Sprintf("UPDATE %s SET description_route = jsonb_set(description_route::jsonb, '{photos}', $1::jsonb) WHERE id = $2", tourTable)
		_, err = tx.Exec(updateQuery, descriptionRouteJSON, tourID)
	}

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
