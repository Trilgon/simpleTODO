package models

type Note struct {
	Id        int     `json:"id" db:"id" validate:"required, gt=0"`
	Email     string  `json:"email" db:"email" validate:"required"`
	Title     string  `json:"title" db:"title" validate:"required"`
	Text      *string `json:"text" db:"text"`
	StartDate string  `json:"start_date" db:"start_date" validate:"required"`
	EndDate   *string `json:"end_date" db:"end_date"`
	IsDone    bool    `json:"is_done" db:"is_done" validate:"required"`
}
