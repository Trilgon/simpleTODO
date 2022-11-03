package dto

type NoteToAdd struct {
	Email string  `json:"email" validate:"required,email"`
	Title string  `json:"title" validate:"required"`
	Text  *string `json:"text"`
}

type NoteToUpdate struct {
	Id    int     `json:"id" validate:"required,gt=0"`
	Title string  `json:"title" validate:"required"`
	Text  *string `json:"text"`
}
