package forms

import "time"

type HelpWithTourForm struct {
	Id       int       `json:"id" db:"id"`
	Name     string    `json:"name" db:"name" binding:"required"`
	Phone    string    `json:"phone" db:"phone" binding:"required"`
	Place    string    `json:"place" db:"place" binding:"required"`
	WhenDate time.Time `json:"when_date" db:"date" binding:"required"`
}
