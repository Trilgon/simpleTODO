package service

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"simpleTODO/internal/models"
	"simpleTODO/internal/models/dto"
	"strconv"
)

func (s *TodoServer) SignUp(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		if len(body) == 0 {
			msg := "body is empty"
			logrus.Errorf("SignUp error, %s. Error: %s", msg, err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
			return
		}
		msg := "failed to read body"
		logrus.Errorf("SignUp error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	user := models.User{}
	err = jsoniter.Unmarshal(body, &user)
	if err != nil {
		msg := "failed to unmarshal body"
		logrus.Errorf("SignUp error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	err = s.rep.SignUp(user.Email, user.Password)
	if err != nil {
		msg := "failed sign up"
		logrus.Errorf("SignUp error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	ctx.Status(http.StatusOK)
	logrus.Infof("SignUp handler called. User was successfully signed up")
}

func (s *TodoServer) SignIn(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		if len(body) == 0 {
			msg := "body is empty"
			logrus.Errorf("SignIn error, %s. Error: %s", msg, err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
			return
		}
		msg := "failed to read body"
		logrus.Errorf("SignIn error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	user := models.User{}
	err = jsoniter.Unmarshal(body, &user)
	if err != nil {
		msg := "failed to unmarshal body"
		logrus.Errorf("SignIn error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
		return
	}
	err = s.rep.SignIn(user.Email, user.Password)
	if err != nil {
		msg := "failed to sign in"
		logrus.Errorf("SignIn error, %s. Errror: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	ctx.Status(http.StatusOK)
	logrus.Infof("SignIn handler called. User was successfully signed in")
}

func (s *TodoServer) SignOut(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		if len(body) == 0 {
			msg := "body is empty"
			logrus.Errorf("SignOut error, %s. Error: %s", msg, err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
			return
		}
		msg := "failed to read body"
		logrus.Errorf("SignOut error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	user := struct {
		Email string `json:"email"`
	}{}
	err = jsoniter.Unmarshal(body, &user)
	if err != nil {
		msg := "failed to unmarshal body"
		logrus.Errorf("SignOut error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
		return
	}
	err = s.rep.SignOut(user.Email)
	if err != nil {
		msg := "failed to sign in"
		logrus.Errorf("SignOut error, %s. Errror: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	ctx.Status(http.StatusOK)
	logrus.Infof("SignOut handler called. User was successfully signed out")
}

func (s *TodoServer) GetById(ctx *gin.Context) {
	key := ctx.Request.URL.Query().Get("id")
	if key == "" {
		msg := "id in query required"
		logrus.Errorf("GetById error, %s.", msg)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}
	id, err := strconv.Atoi(key)
	if err != nil {
		msg := "failed to parse int id from query"
		logrus.Errorf("GetById error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": msg})
		return
	}

	note, err := s.rep.GetById(id)
	if err != nil {
		msg := "failed to get note by id"
		logrus.Errorf("GetById error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	resp, err := jsoniter.Marshal(&note)
	if err != nil {
		msg := "failed to marshal note to response body"
		logrus.Errorf("GetById error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	_, err = ctx.Writer.Write(resp)
	if err != nil {
		msg := "failed to write response body"
		logrus.Errorf("GetById error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	logrus.Infof("GetById handler called. Note successfully got")
}

func (s *TodoServer) GetByEmail(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		if len(body) == 0 {
			msg := "body is empty"
			logrus.Errorf("GetByEmail error, %s. Error: %s", msg, err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
			return
		}
		msg := "failed to read body"
		logrus.Errorf("GetByEmail error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	user := struct {
		Email string `json:"email"`
	}{}
	err = jsoniter.Unmarshal(body, &user)
	if err != nil {
		msg := "failed to unmarshal body"
		logrus.Errorf("GetByEmail error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	notes, err := s.rep.GetByEmail(user.Email)
	if err != nil {
		msg := "failed to get notes by email"
		logrus.Errorf("GetByEmail error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	resp, err := jsoniter.Marshal(&notes)
	if err != nil {
		msg := "failed to marshal notes"
		logrus.Errorf("GetByEmail error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	_, err = ctx.Writer.Write(resp)
	if err != nil {
		msg := "failed to write response body"
		logrus.Errorf("GetByEmail error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	logrus.Infof("GetByEmail handler called. Notes successfully got")
}

func (s *TodoServer) SearchByText(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		if len(body) == 0 {
			msg := "body is empty"
			logrus.Errorf("SearchByText error, %s. Error: %s", msg, err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
		}
		msg := "failed to read body"
		logrus.Errorf("SearchByText error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	text := struct {
		Data string `json:"data"`
	}{}
	err = jsoniter.Unmarshal(body, &text)
	if err != nil {
		msg := "failed to unmarshal body"
		logrus.Errorf("SearchByText error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, err)
		return
	}

	notes, err := s.rep.SearchByText(text.Data)
	if err != nil {
		msg := "failed to search notes"
		logrus.Errorf("SearchByText error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	resp, err := jsoniter.Marshal(notes)
	if err != nil {
		msg := "failed to marshal notes"
		logrus.Errorf("SearchByText error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	_, err = ctx.Writer.Write(resp)
	if err != nil {
		msg := "failed to write response body"
		logrus.Errorf("SearchByText error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	logrus.Infof("SearchByText handler called. Notes successfully got")
}

func (s *TodoServer) AddNote(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		if len(body) == 0 {
			msg := "body is empty"
			logrus.Errorf("AddNote error, %s. Error: %s", msg, err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
			return
		}
		msg := "failed to read body"
		logrus.Errorf("AddNote error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	note := dto.NoteToAdd{}
	err = jsoniter.Unmarshal(body, &note)
	if err != nil {
		msg := "failed to unmarshal body"
		logrus.Errorf("AddNote error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
		return
	}

	id, err := s.rep.AddNote(note.Email, note.Title, note.Text)
	if err != nil {
		msg := "failed to add note"
		logrus.Errorf("AddNote error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}

	resp, err := jsoniter.Marshal(&id)
	if err != nil {
		msg := "failed to marshal response"
		logrus.Errorf("AddNote error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	_, err = ctx.Writer.Write(resp)
	if err != nil {
		msg := "failed to write response body"
		logrus.Errorf("AddNote error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	logrus.Infof("AddNote handler called. Note successfully added")
}

func (s *TodoServer) DeleteNotes(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		if len(body) == 0 {
			msg := "body is empty"
			logrus.Errorf("DeleteNotes error, %s. Error: %s", msg, err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
			return
		}
		msg := "failed to read body"
		logrus.Errorf("DeleteNotes error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	id := make([]int, 0)
	err = jsoniter.Unmarshal(body, &id)
	if err != nil {
		msg := "failed to unmarshal body"
		logrus.Errorf("DeleteNotes error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
		return
	}

	err = s.rep.DeleteNotes(id)
	if err != nil {
		msg := "failed to delete notes"
		logrus.Errorf("DeleteNotes error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	ctx.Status(http.StatusOK)
	logrus.Infof("DeleteNotes handler called. Notes successfully deleted")
}

func (s *TodoServer) UpdateNote(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		if len(body) == 0 {
			msg := "body is empty"
			logrus.Errorf("UpdateNote error, %s. Error: %s", msg, err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
			return
		}
		msg := "failed to read body"
		logrus.Errorf("UpdateNote error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	note := dto.NoteToUpdate{}
	err = jsoniter.Unmarshal(body, &note)
	if err != nil {
		msg := "failed to unmarshal body"
		logrus.Errorf("UpdateNote error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
		return
	}

	err = s.rep.UpdateNote(note.Id, note.Title, note.Text)
	if err != nil {
		msg := "failed to update note"
		logrus.Errorf("UpdateNote error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	ctx.Status(http.StatusOK)
	logrus.Infof("UpdateNote handler called. Note successfully updated")
}

func (s *TodoServer) MarkNote(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		if len(body) == 0 {
			msg := "body is empty"
			logrus.Errorf("MarkNote error, %s. Error: %s", msg, err)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
			return
		}
		msg := "failed to read body"
		logrus.Errorf("MarkNote error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	mark := struct {
		Id    int  `json:"id" db:"id" validate:"required, gt=0"`
		State bool `json:"state" db:"state" validate:"required"`
	}{}
	err = jsoniter.Unmarshal(body, &mark)
	if err != nil {
		msg := "failed to unmarshal body"
		logrus.Errorf("MarkNote error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": msg})
		return
	}

	err = s.rep.MarkNote(mark.Id, mark.State)
	if err != nil {
		msg := "failed to mark note"
		logrus.Errorf("MarkNote error, %s. Error: %s", msg, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	ctx.Status(http.StatusOK)
	logrus.Infof("MarkNote handler called. Note successfully marked")
}
