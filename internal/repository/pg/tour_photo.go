package pg

import (
	"fmt"
	"mime/multipart"
	"silkroad/m/internal/aws"
	"silkroad/m/internal/domain/tour"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type TourPhotoPostgres struct {
	db *sqlx.DB
}

func NewTourPhoto(db *sqlx.DB) *TourPhotoPostgres {
	return &TourPhotoPostgres{db: db}
}

func (r *TourPhotoPostgres) Create(photo tour.TourPhotoInput, photoUrl string) (int, error) {
	var id int
	query := fmt.Sprintf(`
		INSERT INTO %s (tour_id, photo_url, photo_type, title, description, alt_text, display_order, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`, tourPhotosTable)

	err := r.db.QueryRow(query, photo.TourID, photoUrl, photo.PhotoType, photo.Title,
		photo.Description, photo.AltText, photo.DisplayOrder, true).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create tour photo: %w", err)
	}

	return id, nil
}

func (r *TourPhotoPostgres) GetByID(id int) (tour.TourPhoto, error) {
	var photo tour.TourPhoto
	query := fmt.Sprintf(`
		SELECT id, tour_id, photo_url, photo_type, title, description, alt_text, 
		       display_order, is_active, created_at, updated_at
		FROM %s 
		WHERE id = $1 AND is_active = true`, tourPhotosTable)

	err := r.db.Get(&photo, query, id)
	if err != nil {
		return tour.TourPhoto{}, fmt.Errorf("failed to get tour photo: %w", err)
	}

	return photo, nil
}

func (r *TourPhotoPostgres) GetByTourID(tourID int) (*tour.TourPhotosGrouped, error) {
	query := fmt.Sprintf(`
		SELECT id, tour_id, photo_url, photo_type, title, description, alt_text, 
		       display_order, is_active, created_at, updated_at
		FROM %s 
		WHERE tour_id = $1 AND is_active = true
		ORDER BY photo_type, display_order, created_at`, tourPhotosTable)

	var photos []tour.TourPhoto
	err := r.db.Select(&photos, query, tourID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tour photos: %w", err)
	}

	grouped := &tour.TourPhotosGrouped{
		Preview: make([]tour.TourPhoto, 0),
		Gallery: make([]tour.TourPhoto, 0),
		Route:   make([]tour.TourPhoto, 0),
		Booking: make([]tour.TourPhoto, 0),
	}

	for _, photo := range photos {
		switch photo.PhotoType {
		case tour.PhotoTypePreview:
			grouped.Preview = append(grouped.Preview, photo)
		case tour.PhotoTypeGallery:
			grouped.Gallery = append(grouped.Gallery, photo)
		case tour.PhotoTypeRoute:
			grouped.Route = append(grouped.Route, photo)
		case tour.PhotoTypeBooking:
			grouped.Booking = append(grouped.Booking, photo)
		}
	}

	return grouped, nil
}

func (r *TourPhotoPostgres) GetByFilter(filter tour.TourPhotoFilter) ([]tour.TourPhoto, int, error) {
	var conditions []string
	var args []interface{}
	argCount := 0

	baseQuery := fmt.Sprintf(`
		SELECT id, tour_id, photo_url, photo_type, title, description, alt_text, 
		       display_order, is_active, created_at, updated_at
		FROM %s`, tourPhotosTable)

	if filter.TourID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("tour_id = $%d", argCount))
		args = append(args, *filter.TourID)
	}

	if len(filter.PhotoType) > 0 {
		argCount++
		photoTypes := make([]string, len(filter.PhotoType))
		for i, pt := range filter.PhotoType {
			photoTypes[i] = string(pt)
		}
		conditions = append(conditions, fmt.Sprintf("photo_type = ANY($%d)", argCount))
		args = append(args, "{"+strings.Join(photoTypes, ",")+"}")
	}

	if filter.IsActive != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argCount))
		args = append(args, *filter.IsActive)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s%s", tourPhotosTable, whereClause)
	var total int
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tour photos: %w", err)
	}

	query := baseQuery + whereClause + " ORDER BY display_order, created_at"

	if filter.Limit > 0 {
		argCount++
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, filter.Limit)
	}

	if filter.Offset > 0 {
		argCount++
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, filter.Offset)
	}

	var photos []tour.TourPhoto
	err = r.db.Select(&photos, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get tour photos: %w", err)
	}

	return photos, total, nil
}

func (r *TourPhotoPostgres) Update(id int, photo tour.TourPhotoInput) error {
	query := fmt.Sprintf(`
		UPDATE %s 
		SET photo_type = $1, title = $2, description = $3, alt_text = $4, 
		    display_order = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6`, tourPhotosTable)

	_, err := r.db.Exec(query, photo.PhotoType, photo.Title, photo.Description,
		photo.AltText, photo.DisplayOrder, id)
	if err != nil {
		return fmt.Errorf("failed to update tour photo: %w", err)
	}

	return nil
}

func (r *TourPhotoPostgres) Delete(id int) error {
	query := fmt.Sprintf("UPDATE %s SET is_active = false WHERE id = $1", tourPhotosTable)
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tour photo: %w", err)
	}
	return nil
}

func (r *TourPhotoPostgres) DeleteByTourID(tourID int) error {
	query := fmt.Sprintf("UPDATE %s SET is_active = false WHERE tour_id = $1", tourPhotosTable)
	_, err := r.db.Exec(query, tourID)
	if err != nil {
		return fmt.Errorf("failed to delete tour photos: %w", err)
	}
	return nil
}

func (r *TourPhotoPostgres) UploadPhotos(tourID int, files []*multipart.FileHeader, photoType tour.TourPhotoType, metadata tour.TourPhotoInput) error {
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
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	for i, file := range files {
		f, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}

		key := fmt.Sprintf("tour_photos/%s/%d_%d_%s", photoType, tourID, time.Now().Unix(), file.Filename)

		fileUrl, err := aws.UploadPhotoToS3(aws.GetBucketName(), key, f)
		if err != nil {
			_ = f.Close()
			return fmt.Errorf("failed to upload photo to S3: %w", err)
		}

		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close file: %w", err)
		}

		insertQuery := fmt.Sprintf(`
			INSERT INTO %s (tour_id, photo_url, photo_type, title, description, alt_text, display_order, is_active)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, tourPhotosTable)

		_, err = tx.Exec(insertQuery, tourID, fileUrl, photoType, metadata.Title,
			metadata.Description, metadata.AltText, metadata.DisplayOrder+i, true)
		if err != nil {
			return fmt.Errorf("failed to save photo to database: %w", err)
		}
	}

	return tx.Commit()
}
