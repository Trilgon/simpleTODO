package postgresql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"simpleTODO/internal/models"
	"simpleTODO/pkg/db/postgresql"
	"time"
)

type TodoRepository struct {
	db *sqlx.DB
}

func NewTodoRepository() (*TodoRepository, error) {
	db, err := postgresql.GetDB()
	if err != nil {
		return nil, err
	}
	return &TodoRepository{db: db}, nil
}

func (p *TodoRepository) checkExists(email string) (bool, error) {
	var isExists bool
	err := p.db.Get(&isExists, checkUser, email)
	if err != nil {
		return false, fmt.Errorf("failed to check if already user exists. Error: %s", err)
	}
	return isExists, nil
}

func (p *TodoRepository) SignUp(email, encPassword string) error {
	isExists, err := p.checkExists(email)
	if err != nil {
		return err
	}
	if isExists {
		return fmt.Errorf("user already exists")
	}
	_, err = p.db.Exec(newUser, email, encPassword)
	if err != nil {
		return fmt.Errorf("failed to add new user. Error: %s", err)
	}
	return nil
}

func (p *TodoRepository) SignIn(email, encPassword string) error {
	var id int
	err := p.db.QueryRow(authorizeUser, email, encPassword).Scan(&id)
	if err != nil {
		return fmt.Errorf("failed to check if user exists. Error: %s", err)
	}
	if id < 1 {
		return fmt.Errorf("user doesn't exists")
	}
	_, err = p.db.Exec(signIn, id)
	if err != nil {
		return fmt.Errorf("failed to mark user as logged. Error: %s", err)
	}
	return nil
}

func (p *TodoRepository) SignOut(email string) error {
	isExists, err := p.checkExists(email)
	if err != nil {
		return err
	}
	if !isExists {
		return fmt.Errorf("user doesn't exists")
	}
	_, err = p.db.Exec(signOut, email)
	if err != nil {
		return fmt.Errorf("failed to sign out user. Error: %s", err)
	}
	return nil
}

func (p *TodoRepository) AddNote(email, title, text string) (int, error) {
	currentDt := time.Now()
	var id int
	err := p.db.QueryRow(addNote, email, title, text, currentDt.Format("2006-01-02 15:04:05")).
		Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("failed to add note. Error: %s", err)
	}
	return id, nil
}

func (p *TodoRepository) DeleteNote(id int) error {
	_, err := p.db.Exec(deleteNote, id)
	if err != nil {
		return fmt.Errorf("failed to delete note. Error: %s", err)
	}
	return nil
}

func (p *TodoRepository) UpdateNote(id int, title, text string) error {
	_, err := p.db.Exec(updateNote, title, text, id)
	if err != nil {
		return fmt.Errorf("failed to update note. Error: %s", err)
	}
	return nil
}

func (p *TodoRepository) MarkNote(id int, state bool) error {
	_, err := p.db.Exec(markNote, state, id)
	if err != nil {
		return fmt.Errorf("failed to mark note. Error: %s", err)
	}
	return nil
}

func (p *TodoRepository) GetById(id int) (*models.Note, error) {
	note := models.Note{}
	err := p.db.Get(&note, getById)
	if err != nil {
		return nil, fmt.Errorf("failed to get note with id = %d. Error: %s", id, err)
	}
	return &note, nil
}

func (p *TodoRepository) GetByEmail(email string) ([]models.Note, error) {
	notes := make([]models.Note, 0)
	err := p.db.Select(&notes, getByEmail, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get notes where email = %s. Error: %s", email, err)
	}
	return notes, nil
}

func (p *TodoRepository) SearchByText(text string) ([]models.Note, error) {
	//TODO implement me
	panic("implement me")
}
