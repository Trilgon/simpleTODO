package dto

type NoteToAdd struct {
	Email string  `json:"email" db:"email" validate:"required"`
	Title string  `json:"title" db:"title" validate:"required"`
	Text  *string `json:"text" db:"text"`
}

type NoteToUpdate struct {
	Id    int    `json:"id" db:"id" validate:"required, gt=0"`
	Title string `json:"title" db:"title" validate:"required"`
	Text  string `json:"text" db:"text"`
}
