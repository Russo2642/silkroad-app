package tour

import "time"

type TourEditor struct {
	Id       int       `json:"id" db:"id"`
	Name     string    `json:"name" db:"name" binding:"required"`
	Phone    string    `json:"phone" db:"phone" binding:"required"`
	Email    string    `json:"email" db:"email" binding:"required"`
	TourDate time.Time `json:"tour_date" db:"tour_date" binding:"required"`
	Activity []string  `json:"activity" db:"activity" binding:"required"`
	Location []string  `json:"location" db:"location" binding:"required"`
}
