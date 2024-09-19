package forms

type ContactForm struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name" binding:"required"`
	Phone       string `json:"phone" db:"phone" binding:"required"`
	Email       string `json:"email" db:"email" binding:"required"`
	Description string `json:"description" db:"description" binding:"required"`
	TourID      *int   `json:"tour_id" db:"tour_id"`
}
