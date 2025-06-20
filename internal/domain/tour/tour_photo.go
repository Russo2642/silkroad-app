package tour

import "time"

type TourPhotoType string

const (
	PhotoTypePreview TourPhotoType = "preview"
	PhotoTypeGallery TourPhotoType = "gallery"
	PhotoTypeRoute   TourPhotoType = "route"
	PhotoTypeBooking TourPhotoType = "booking"
)

type TourPhoto struct {
	ID           int           `json:"id" db:"id"`
	TourID       int           `json:"tour_id" db:"tour_id"`
	PhotoUrl     string        `json:"photo_url" db:"photo_url" binding:"required"`
	PhotoType    TourPhotoType `json:"photo_type" db:"photo_type" binding:"required"`
	Title        *string       `json:"title,omitempty" db:"title"`
	Description  *string       `json:"description,omitempty" db:"description"`
	AltText      *string       `json:"alt_text,omitempty" db:"alt_text"`
	DisplayOrder int           `json:"display_order" db:"display_order"`
	IsActive     bool          `json:"is_active" db:"is_active"`
	CreatedAt    time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at" db:"updated_at"`
}

type TourPhotoInput struct {
	TourID       int           `json:"tour_id" binding:"required"`
	PhotoType    TourPhotoType `json:"photo_type" binding:"required"`
	Title        *string       `json:"title,omitempty"`
	Description  *string       `json:"description,omitempty"`
	AltText      *string       `json:"alt_text,omitempty"`
	DisplayOrder int           `json:"display_order"`
}

type TourPhotoFilter struct {
	TourID    *int            `json:"tour_id" form:"tour_id"`
	PhotoType []TourPhotoType `json:"photo_type" form:"photo_type"`
	IsActive  *bool           `json:"is_active" form:"is_active"`
	Limit     int             `json:"limit" form:"limit"`
	Offset    int             `json:"offset" form:"offset"`
}

type TourPhotosGrouped struct {
	Preview []TourPhoto `json:"preview"`
	Gallery []TourPhoto `json:"gallery"`
	Route   []TourPhoto `json:"route"`
	Booking []TourPhoto `json:"booking"`
}

func IsValidPhotoType(t TourPhotoType) bool {
	switch t {
	case PhotoTypePreview, PhotoTypeGallery, PhotoTypeRoute, PhotoTypeBooking:
		return true
	}
	return false
}
