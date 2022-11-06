package db

import (
	"todoServer/internal/models"
)

type TodoRepository interface {
	//SignUp(email, encPassword string) error
	//SignIn(email, encPassword string) error
	//SignOut(email string) error
	AddNote(email, title string, text *string) (int, error)
	DeleteNotes(email string, id []int) error
	UpdateNote(id int, title string, text *string, email string) error
	MarkNote(id int, state bool, email string) error
	GetById(id int, email string) (*models.Note, error)
	GetByEmail(email string) ([]models.Note, error)
	SearchByText(email string, text string) ([]models.Note, error)
}
