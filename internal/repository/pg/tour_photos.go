package pg

import (
	"database/sql"
	"fmt"
	"mime/multipart"
	"silkroad/m/internal/aws"
	"time"
)

func (r *TourPostgres) AddPhotos(tourID int, files []*multipart.FileHeader, photoType string) error {
	var exists bool
	checkQuery := fmt.Sprintf(`
		SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)`,
		tourTable)

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

	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		key := fmt.Sprintf("tour_photos/%s/%d_%s", photoType, time.Now().Unix(), file.Filename)
		fileUrl, err := aws.UploadPhotoToS3("gosilkroadbucket", key, f)
		if err != nil {
			_ = f.Close()
			return fmt.Errorf("failed to upload photo: %w", err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}

		insertQuery := fmt.Sprintf("INSERT INTO %s (tour_id, photo_type, photo_url) VALUES ($1, $2, $3)", tourPhotosTable)
		_, err = tx.Exec(insertQuery, tourID, photoType, fileUrl)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *TourPostgres) getPhotosByTourID(tourID int) ([]string, []string, string, string, error) {
	var (
		galleryPhotos []string
		routePhotos   []string
		previewPhoto  string
		bookTourPhoto string
	)

	query := fmt.Sprintf(`
        SELECT photo_type, photo_url
        FROM %s
        WHERE tour_id = $1
    `, tourPhotosTable)

	rows, err := r.db.Query(query, tourID)
	if err != nil {
		return nil, nil, "", "", err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var photoType string
		var photoUrl string

		if err := rows.Scan(&photoType, &photoUrl); err != nil {
			return nil, nil, "", "", err
		}

		switch photoType {
		case "gallery":
			galleryPhotos = append(galleryPhotos, photoUrl)
		case "route":
			routePhotos = append(routePhotos, photoUrl)
		case "preview":
			previewPhoto = photoUrl
		case "book":
			bookTourPhoto = photoUrl
		}
	}

	if err = rows.Err(); err != nil {
		return nil, nil, "", "", err
	}

	return galleryPhotos, routePhotos, previewPhoto, bookTourPhoto, nil
}

//func (r *TourPostgres) AddPhotos(tourID int, files []*multipart.FileHeader, updateField string) error {
//	var exists bool
//	checkQuery := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE id = $1)", tourTable)
//	err := r.db.QueryRow(checkQuery, tourID).Scan(&exists)
//	if err != nil {
//		return fmt.Errorf("failed to check tour existence: %w", err)
//	}
//
//	if !exists {
//		return fmt.Errorf("tour with id %d does not exist", tourID)
//	}
//
//	tx, err := r.db.Begin()
//	if err != nil {
//		return err
//	}
//
//	defer func() {
//		if err != nil {
//			_ = tx.Rollback()
//		}
//	}()
//
//	var photoUrls []string
//
//	for _, file := range files {
//		f, err := file.Open()
//		if err != nil {
//			return fmt.Errorf("failed to open file: %w", err)
//		}
//
//		key := fmt.Sprintf("tour_photos/%d_%s", time.Now().Unix(), file.Filename)
//
//		fileUrl, err := aws.UploadPhotoToS3("gosilkroadbucket", key, f)
//		if err != nil {
//			_ = f.Close()
//			return fmt.Errorf("failed to upload photo: %w", err)
//		}
//
//		photoUrls = append(photoUrls, fileUrl)
//
//		if err := f.Close(); err != nil {
//			return fmt.Errorf("failed to close file: %w", err)
//		}
//	}
//
//	var updateQuery string
//	if updateField == "gallery" {
//		updateQuery = fmt.Sprintf("UPDATE %s SET photos = array_cat(COALESCE(photos, ARRAY[]::text[]), $1) WHERE id = $2", tourTable)
//		_, err = tx.Exec(updateQuery, pq.Array(photoUrls), tourID)
//	} else if updateField == "route" {
//		descriptionRouteJSON, err := json.Marshal(photoUrls)
//		if err != nil {
//			return fmt.Errorf("failed to marshal photo URLs: %w", err)
//		}
//		updateQuery = fmt.Sprintf("UPDATE %s SET description_route = jsonb_set(description_route::jsonb, '{photos}', $1::jsonb) WHERE id = $2", tourTable)
//		_, err = tx.Exec(updateQuery, descriptionRouteJSON, tourID)
//	}
//
//	if err != nil {
//		_ = tx.Rollback()
//		return err
//	}
//
//	return tx.Commit()
//}
