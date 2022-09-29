package db

import "simpleTODO/internal/models"

type Repository interface {
	SignUp(email, encPassword string) error
	SignIn(email, encPassword string) error
	SignOut(email string) error
	AddNote(email, title, text string) (int, error)
	DeleteNote(id int) error
	UpdateNote(id int, title, text string) error
	MarkNote(id int, state bool) error
	GetById(id int) (*models.Note, error)
	GetByEmail(email string) ([]models.Note, error)
	SearchByText(text string) ([]models.Note, error)
}
