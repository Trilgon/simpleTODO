package postgresql

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"time"
	"todoServer/internal/models"
)

type TodoRepository struct {
	db *sqlx.DB
}

func NewTodoRepository() (*TodoRepository, error) {
	db, err := GetDB()
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

//func (p *TodoRepository) SignUp(email, encPassword string) error {
//	isExists, err := p.checkExists(email)
//	if err != nil {
//		return err
//	}
//	if isExists {
//		return fmt.Errorf("user already exists")
//	}
//	_, err = p.db.Exec(newUser, email, encPassword)
//	if err != nil {
//		return fmt.Errorf("failed to add new user. Error: %s", err)
//	}
//	return nil
//}

//func (p *TodoRepository) SignIn(email, encPassword string) error {
//	var id int
//	err := p.db.QueryRow(authorizeUser, email, encPassword).Scan(&id)
//	if err != nil {
//		return fmt.Errorf("failed to check if user exists. Error: %s", err)
//	}
//	if id < 1 {
//		return fmt.Errorf("user doesn't exists")
//	}
//	_, err = p.db.Exec(signIn, email)
//	if err != nil {
//		return fmt.Errorf("failed to mark user as logged. Error: %s", err)
//	}
//	return nil
//}

//func (p *TodoRepository) SignOut(email string) error {
//	isExists, err := p.checkExists(email)
//	if err != nil {
//		return err
//	}
//	if !isExists {
//		return fmt.Errorf("user doesn't exists")
//	}
//	_, err = p.db.Exec(signOut, email)
//	if err != nil {
//		return fmt.Errorf("failed to sign out user. Error: %s", err)
//	}
//	return nil
//}

func (p *TodoRepository) AddNote(email, title string, text *string) (int, error) {
	currentDt := time.Now()
	var id int
	var err error
	if text != nil {
		err = p.db.QueryRow(addNote, email, title, *text, currentDt.Format("2006-01-02 15:04:05")).
			Scan(&id)
	} else {
		err = p.db.QueryRow(addNote, email, title, nil, currentDt.Format("2006-01-02 15:04:05")).
			Scan(&id)
	}
	if err != nil {
		return -1, fmt.Errorf("failed to add note. Error: %s", err)
	}
	return id, nil
}

func (p *TodoRepository) DeleteNotes(email string, id []int) error {
	res, err := p.db.Exec(deleteNotes, email, pq.Array(id))
	if err != nil {
		return fmt.Errorf("failed to delete notes. Error: %s", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected < 1 {
		return fmt.Errorf("notes with ids = %v and email = %s don't exist", id, email)
	}
	return nil
}

func (p *TodoRepository) UpdateNote(id int, title string, text *string, email string) error {
	var res sql.Result
	var err error
	if text != nil {
		res, err = p.db.Exec(updateNote, title, *text, id, email)
	} else {
		res, err = p.db.Exec(updateNote, title, nil, id, email)
	}
	if err != nil {
		return fmt.Errorf("failed to update note. Error: %s", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected < 1 {
		return fmt.Errorf("note with id = %d doesn't exist", id)
	}
	return nil
}

func (p *TodoRepository) MarkNote(id int, state bool, email string) error {
	var res sql.Result
	var err error
	if state {
		res, err = p.db.Exec(markDone, time.Now().Format("2006-01-02 15:04:05"), id, email)
		if err != nil {
			return fmt.Errorf("failed to mark note. Error: %s", err)
		}
	} else {
		res, err = p.db.Exec(markUndone, id, email)
		if err != nil {
			return fmt.Errorf("failed to mark note. Error: %s", err)
		}
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected < 1 {
		return fmt.Errorf("note with id = %d and email = %s doesn't exist", id, email)
	}
	return nil
}

func (p *TodoRepository) GetById(id int, email string) (*models.Note, error) {
	note := models.Note{}
	err := p.db.Get(&note, getById, id, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("note with id = %d and email = %s doesn't exist", id, email)
		}
		return nil, fmt.Errorf("failed to get note with id = %d. Error: %s", id, err)
	}
	return &note, nil
}

func (p *TodoRepository) GetByEmail(email string) ([]models.Note, error) {
	notes := make([]models.Note, 0)
	err := p.db.Select(&notes, getByEmail, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notes with email = %s don't exist", email)
		}
		return nil, fmt.Errorf("failed to get notes where email = %s. Error: %s", email, err)
	}
	return notes, nil
}

func (p *TodoRepository) SearchByText(email string, text string) ([]models.Note, error) {
	notes := make([]models.Note, 0)
	likeParam := "%" + text + "%"
	err := p.db.Select(&notes, getByText, email, likeParam)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notes with email = %s don't exist", email)
		}
		return nil, fmt.Errorf("failed to find any notes that contains text in title or text. Error: %s", err)
	}
	return notes, nil
}
