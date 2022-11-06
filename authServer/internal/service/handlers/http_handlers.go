package handlers

import (
	db "authServer/internal/db"
	"authServer/internal/models"
	"authServer/internal/service/auth"
	"authServer/internal/utils"
	"crypto/sha512"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type HttpHandlers struct {
	rep       db.AuthRepository
	validator *validator.Validate
	ts        ut.Translator
}

func NewHttpHandlers(rep db.AuthRepository, val *validator.Validate, ts ut.Translator) (*HttpHandlers, error) {
	if rep == nil || val == nil || ts == nil {
		return nil, fmt.Errorf("parameters should not be nil. Passed parameters: "+
			"repository: %v, validator: %v, translator: %v", rep, val, ts)
	}
	httpHandlers := HttpHandlers{
		rep:       rep,
		validator: val,
		ts:        ts,
	}
	return &httpHandlers, nil
}

func (h *HttpHandlers) catchErrGin(ctx *gin.Context, code int, msg string, err error) {
	if err == nil {
		logrus.Errorf("%s.", msg)
		ctx.AbortWithStatusJSON(code, gin.H{"error": fmt.Sprintf("%s.", msg)})
		return
	}
	logrus.Errorf("%s. Error: %s", msg, err)
	ctx.AbortWithStatusJSON(code, gin.H{"error": fmt.Sprintf("%s. %s", msg, err)})
}

func (h *HttpHandlers) encryptPas(password string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(password)))
}

func (h *HttpHandlers) SignUp(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		if len(body) == 0 {
			h.catchErrGin(ctx, http.StatusUnprocessableEntity, "SignUp error, body is empty", err)
		}
		h.catchErrGin(ctx, http.StatusUnprocessableEntity, "SignUp error, failed to read request body", err)
		return
	}
	var user models.User
	err = jsoniter.Unmarshal(body, &user)
	if err != nil {
		h.catchErrGin(ctx, http.StatusUnprocessableEntity, "SignUp error, failed to unmarshal request body", err)
		return
	}

	err = h.validator.Struct(user)
	if err != nil {
		h.catchErrGin(ctx, http.StatusBadRequest,
			"SignUp error failed to validate user model", utils.TranslateError(err, h.ts))
		return
	}
	err = h.rep.SignUp(user.Email, h.encryptPas(user.Password))
	if err != nil {
		h.catchErrGin(ctx, http.StatusInternalServerError, "SignUp error", err)
		return
	}
	token, err := auth.GenerateJwt(user.Email)
	if err != nil {
		h.catchErrGin(ctx, http.StatusInternalServerError, "SignUp error", err)
		return
	}

	resp, err := jsoniter.Marshal(token)
	if err != nil {
		h.catchErrGin(ctx, http.StatusInternalServerError, "SignUp error, failed to marshal response body", err)
		return
	}
	_, err = ctx.Writer.Write(resp)
	if err != nil {
		h.catchErrGin(ctx, http.StatusInternalServerError, "SignUp error, failed to write response body", err)
		return
	}
	ctx.Status(http.StatusOK)
}

func (h *HttpHandlers) SignIn(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.catchErrGin(ctx, http.StatusUnprocessableEntity, "SignIn error, failed to read request body", err)
		return
	}
	var user models.User
	err = jsoniter.Unmarshal(body, &user)
	if err != nil {
		h.catchErrGin(ctx, http.StatusUnprocessableEntity, "SignIn error, failed to unmarshal request body", err)
		return
	}

	err = h.validator.Struct(user)
	if err != nil {
		h.catchErrGin(ctx, http.StatusBadRequest,
			"SignIn error, failed to validate user model", utils.TranslateError(err, h.ts))
		return
	}
	err = h.rep.SignIn(user.Email, h.encryptPas(user.Password))
	if err != nil {
		h.catchErrGin(ctx, http.StatusUnauthorized, "SignIn error", err)
		return
	}
	token, err := auth.GenerateJwt(user.Email)
	if err != nil {
		h.catchErrGin(ctx, http.StatusInternalServerError, "SignIn error", err)
		return
	}

	resp, err := jsoniter.Marshal(token)
	if err != nil {
		h.catchErrGin(ctx, http.StatusInternalServerError, "SignIn error, failed to marshal response body", err)
		return
	}
	_, err = ctx.Writer.Write(resp)
	if err != nil {
		h.catchErrGin(ctx, http.StatusInternalServerError, "SignIn error, failed to write response body", err)
		return
	}
	ctx.Status(http.StatusOK)
}
