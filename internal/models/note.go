package models

type Note struct {
	Id        int     `json:"id" db:"id"`
	Email     string  `json:"email" db:"email"`
	Title     string  `json:"title" db:"title"`
	Text      *string `json:"text" db:"text"`
	StartDate string  `json:"start_date" db:"start_date"`
	EndDate   *string `json:"end_date" db:"end_date"`
	IsDone    bool    `json:"is_done" db:"is_done"`
}
