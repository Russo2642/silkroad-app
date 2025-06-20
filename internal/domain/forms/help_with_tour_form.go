package forms

import "time"

type HelpWithTourForm struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required,min=2,max=100"`
	Phone     string    `json:"phone" db:"phone" binding:"required,min=10,max=20"`
	Country   string    `json:"country" db:"country" binding:"required,min=2,max=100"`
	WhenDate  time.Time `json:"when_date" db:"when_date" binding:"required"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
