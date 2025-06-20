package forms

import (
	"time"
)

type ContactForm struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required,min=2,max=100"`
	Phone     string    `json:"phone" db:"phone" binding:"required,min=10,max=20"`
	Email     string    `json:"email" db:"email" binding:"required,email,max=100"`
	TourID    *int      `json:"tour_id" db:"tour_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
