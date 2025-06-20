package country

import "time"

type Country struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required"`
	Code      string    `json:"code" db:"code" binding:"required"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CountryFilter struct {
	IsActive *bool `json:"is_active" form:"is_active"`
}
