package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
	"strings"
	"todoServer/internal/models/dto"
	"todoServer/internal/service/auth"
	"todoServer/internal/utils"
)

const emailHeader = "X-Email"

func (s *TodoServer) catchErrGin(ctx *gin.Context, code int, msg string, err error) {
	if err == nil {
		logrus.Errorf("%s.", msg)
		ctx.AbortWithStatusJSON(code, gin.H{"error": fmt.Sprintf("%s.", msg)})
		return
	}
	logrus.Errorf("%s. Error: %s", msg, err)
	ctx.AbortWithStatusJSON(code, gin.H{"error": fmt.Sprintf("%s. %s", msg, err)})
}

func (s *TodoServer) AuthMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		s.catchErrGin(ctx, http.StatusUnauthorized, "Auth error, no auth data in request header", nil)
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		s.catchErrGin(ctx, http.StatusUnauthorized,
			"Auth error, header Authorization should looks like \"Bearer tokenStr\"", nil)
		return
	}
	publicKey, err := auth.GetRSAPublicKey()
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "Auth error", err)
	}
	claimsMap, err := auth.ValidateToken(parts[1], publicKey)
	if err != nil {
		s.catchErrGin(ctx, http.StatusUnauthorized, "Auth error", err)
		return
	}

	email := claimsMap["email"]
	err = s.validator.Var(email, "required,email")
	if err != nil {
		s.catchErrGin(ctx, http.StatusUnauthorized,
			"Auth error, email must be set in jwt claims", utils.TranslateError(err, s.ts))
		return
	}
	ctx.Request.Header.Set(emailHeader, email.(string))
	ctx.Next()
}

//func (s *TodoServer) SignUp(ctx *gin.Context) {
//	body, err := io.ReadAll(ctx.Request.Body)
//	if err != nil {
//		if len(body) == 0 {
//			s.catchErrGin(ctx, http.StatusUnprocessableEntity, "SignUp error, body is empty", err)
//			return
//		}
//		s.catchErrGin(ctx, http.StatusInternalServerError, "SignUp error, failed to read body", err)
//		return
//	}
//	user := models.User{}
//	err = jsoniter.Unmarshal(body, &user)
//	if err != nil {
//		s.catchErrGin(ctx, http.StatusInternalServerError, "SignUp error, failed to unmarshal body", err)
//		return
//	}
//
//	err = s.validator.Struct(user)
//	if err != nil {
//		s.catchErrGin(ctx, http.StatusBadRequest,
//			"SignUp error, failed to validate body", utils.TranslateError(err, s.ts))
//		return
//	}
//	err = s.rep.SignUp(user.Email, user.Password)
//	if err != nil {
//		s.catchErrGin(ctx, http.StatusInternalServerError, "SignUp error, failed to sign up", err)
//		return
//	}
//	ctx.Status(http.StatusOK)
//	logrus.Infof("SignUp handler called. User was successfully signed up")
//}

//func (s *TodoServer) SignIn(ctx *gin.Context) {
//	body, err := io.ReadAll(ctx.Request.Body)
//	if err != nil {
//		if len(body) == 0 {
//			s.catchErrGin(ctx, http.StatusUnprocessableEntity, "SignIn error, body is empty", err)
//			return
//		}
//		s.catchErrGin(ctx, http.StatusInternalServerError, "SignIn error, failed to read body", err)
//		return
//	}
//	user := models.User{}
//	err = jsoniter.Unmarshal(body, &user)
//	if err != nil {
//		s.catchErrGin(ctx, http.StatusUnprocessableEntity, "SignIn error, failed to unmarshal body", err)
//		return
//	}
//	err = s.validator.Struct(user)
//	if err != nil {
//		s.catchErrGin(ctx, http.StatusBadRequest,
//			"SignIn error, failed to validate body.", utils.TranslateError(err, s.ts))
//		return
//	}
//	err = s.rep.SignIn(user.Email, user.Password)
//	if err != nil {
//		s.catchErrGin(ctx, http.StatusInternalServerError, "SignIn error, failed to sign in", err)
//		return
//	}
//	ctx.Status(http.StatusOK)
//	logrus.Infof("SignIn handler called. User was successfully signed in")
//}
//
//func (s *TodoServer) SignOut(ctx *gin.Context) {
//	body, err := io.ReadAll(ctx.Request.Body)
//	if err != nil {
//		if len(body) == 0 {
//			s.catchErrGin(ctx, http.StatusUnprocessableEntity, "SignOut error, body is empty", err)
//			return
//		}
//		s.catchErrGin(ctx, http.StatusInternalServerError, "SignOut error, failed to read body", err)
//		return
//	}
//	user := struct {
//		Email string `json:"email"`
//	}{}
//	err = jsoniter.Unmarshal(body, &user)
//	if err != nil {
//		s.catchErrGin(ctx, http.StatusUnprocessableEntity, "SignOut error, failed to unmarshal body", err)
//		return
//	}
//
//	err = s.validator.Struct(user)
//	if err != nil {
//		s.catchErrGin(ctx, http.StatusBadRequest,
//			"SignOut error, failed to validate body", utils.TranslateError(err, s.ts))
//		return
//	}
//	err = s.rep.SignOut(user.Email)
//	if err != nil {
//		s.catchErrGin(ctx, http.StatusInternalServerError, "SignOut error, failed to sign out", err)
//		return
//	}
//	ctx.Status(http.StatusOK)
//	logrus.Infof("SignOut handler called. User with email: %s was successfully signed out", user.Email)
//}

func (s *TodoServer) GetById(ctx *gin.Context) {
	key := ctx.Request.URL.Query().Get("id")
	if key == "" {
		s.catchErrGin(ctx, http.StatusBadRequest, "GetById error", fmt.Errorf("id in query required"))
		return
	}
	id, err := strconv.Atoi(key)
	if err != nil {
		s.catchErrGin(ctx, http.StatusBadRequest, "GetById error, failed to parse int id from query", err)
		return
	}

	err = s.validator.Var(id, "required,gt=0")
	if err != nil {
		s.catchErrGin(ctx, http.StatusBadRequest,
			"GetById error, failed to validate body", utils.TranslateError(err, s.ts))
		return
	}
	email := ctx.GetHeader(emailHeader)
	if email == "" {
		s.catchErrGin(ctx, http.StatusBadRequest, "AddNote error, email must be set", nil)
		return
	}
	note, err := s.rep.GetById(id, email)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "GetById error, failed to get note by id", err)
		return
	}

	resp, err := jsoniter.Marshal(&note)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "GetById error, failed to marshal note to response body", err)
		return
	}
	_, err = ctx.Writer.Write(resp)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "GetById error, failed to write response body", err)
		return
	}
	logrus.Infof("GetById handler called. Note with id %d successfully got", note.Id)
}

func (s *TodoServer) GetByEmail(ctx *gin.Context) {
	email := ctx.GetHeader(emailHeader)
	if email == "" {
		s.catchErrGin(ctx, http.StatusBadRequest, "AddNote error, email must be set", nil)
		return
	}
	notes, err := s.rep.GetByEmail(email)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "GetByEmail error, failed to get notes by email", err)
		return
	}

	resp, err := jsoniter.Marshal(&notes)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "GetByEmail error, failed to marshal note to response body", err)
		return
	}
	_, err = ctx.Writer.Write(resp)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "GetByEmail error, failed to write response body", err)
		return
	}
	ids := make([]int, len(notes), len(notes))
	for i, note := range notes {
		ids[i] = note.Id
	}
	logrus.Infof("GetByEmail handler called. Notes with email: %v and ids: %v successfully got", email, ids)
}

func (s *TodoServer) SearchByText(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "SearchByText error, failed to read body", err)
		return
	}
	text := struct {
		Text string `json:"text"`
	}{}
	err = jsoniter.Unmarshal(body, &text)
	if err != nil {
		s.catchErrGin(ctx, http.StatusUnprocessableEntity, "SearchByText error, failed to unmarshal body", err)
		return
	}

	email := ctx.GetHeader(emailHeader)
	if email == "" {
		s.catchErrGin(ctx, http.StatusBadRequest, "AddNote error, email must be set", nil)
		return
	}
	notes, err := s.rep.SearchByText(email, text.Text)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "SearchByText error, failed to search notes", err)
		return
	}

	resp, err := jsoniter.Marshal(notes)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "SearchByText error, failed to marshal note to response body", err)
		return
	}
	_, err = ctx.Writer.Write(resp)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "SearchByText error, failed to write response body", err)
		return
	}
	ids := make([]int, len(notes), len(notes))
	for i, note := range notes {
		ids[i] = note.Id
	}
	logrus.Infof("SearchByText handler called. Notes with ids: %v successfully got", ids)
}

func (s *TodoServer) AddNote(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "AddNote error, failed to read body", err)
		return
	}
	note := dto.NoteToAdd{}
	err = jsoniter.Unmarshal(body, &note)
	if err != nil {
		s.catchErrGin(ctx, http.StatusUnprocessableEntity, "AddNote error, failed to unmarshal body", err)
		return
	}

	email := ctx.GetHeader(emailHeader)
	if email == "" {
		s.catchErrGin(ctx, http.StatusBadRequest, "AddNote error, email must be set", nil)
		return
	}
	err = s.validator.Struct(note)
	if err != nil {
		s.catchErrGin(ctx, http.StatusBadRequest,
			"AddNote error, failed to validate body", utils.TranslateError(err, s.ts))
		return
	}
	id, err := s.rep.AddNote(email, note.Title, note.Text)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "AddNote error, failed to add note", err)
		return
	}

	resp, err := jsoniter.Marshal(&id)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "AddNote error, failed to marshal note to response body", err)
		return
	}
	_, err = ctx.Writer.Write(resp)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "AddNote error, failed to write response body", err)
		return
	}
	logrus.Infof("AddNote handler called. Note with id: %d successfully added", id)
}

func (s *TodoServer) DeleteNotes(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "DeleteNotes error, failed to read body", err)
		return
	}
	ids := make([]int, 0)
	err = jsoniter.Unmarshal(body, &ids)
	if err != nil {
		s.catchErrGin(ctx, http.StatusUnprocessableEntity, "DeleteNotes error, failed to unmarshal body", err)
		return
	}
	err = s.validator.Var(ids, "required,gt=0")
	if err != nil {
		s.catchErrGin(ctx, http.StatusBadRequest,
			"DeleteNotes error, failed to validate ids", utils.TranslateError(err, s.ts))
		return
	}

	email := ctx.GetHeader(emailHeader)
	if email == "" {
		s.catchErrGin(ctx, http.StatusBadRequest, "AddNote error, email must be set", nil)
		return
	}
	err = s.rep.DeleteNotes(email, ids)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "DeleteNotes error, failed to delete notes", err)
		return
	}
	ctx.Status(http.StatusNoContent)
	logrus.Infof("DeleteNotes handler called. Notes with ids: %v successfully deleted", ids)
}

func (s *TodoServer) UpdateNote(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "UpdateNote error, failed to read body", err)
		return
	}
	note := dto.NoteToUpdate{}
	err = jsoniter.Unmarshal(body, &note)
	if err != nil {
		s.catchErrGin(ctx, http.StatusUnprocessableEntity, "UpdateNote error, failed to unmarshal body", err)
		return
	}

	email := ctx.GetHeader(emailHeader)
	if email == "" {
		s.catchErrGin(ctx, http.StatusBadRequest, "AddNote error, email must be set", nil)
		return
	}
	err = s.rep.UpdateNote(note.Id, note.Title, note.Text, email)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "UpdateNote error, failed to update note", err)
		return
	}
	ctx.Status(http.StatusOK)
	logrus.Infof("UpdateNote handler called. Note with id: %d successfully updated", note.Id)
}

func (s *TodoServer) MarkNote(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "MarkNote error, failed to read body", err)
		return
	}
	mark := struct {
		Id    int  `json:"id" validate:"required, gt=0"`
		State bool `json:"state" validate:"required"`
	}{}
	err = jsoniter.Unmarshal(body, &mark)
	if err != nil {
		s.catchErrGin(ctx, http.StatusUnprocessableEntity, "MarkNote error, failed to unmarshal body", err)
		return
	}

	email := ctx.GetHeader(emailHeader)
	if email == "" {
		s.catchErrGin(ctx, http.StatusBadRequest, "AddNote error, email must be set", nil)
		return
	}
	err = s.rep.MarkNote(mark.Id, mark.State, email)
	if err != nil {
		s.catchErrGin(ctx, http.StatusInternalServerError, "MarkNote error, failed to mark note", err)
		return
	}
	ctx.Status(http.StatusOK)
	logrus.Infof("MarkNote handler called. Note with id: %d successfully marked", mark.Id)
}
