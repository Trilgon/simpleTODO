package service

import (
	"authServer/internal/db/postgresql"
	"authServer/internal/service/handlers"
	"fmt"
	"github.com/gin-gonic/gin"
	eng "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/spf13/viper"
)

type AuthService struct {
	handlers handlers.AuthHandlers
}

func NewAuthService() (*AuthService, error) {
	srv := AuthService{}
	rep, err := postgresql.NewAuthRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to init service. Error: %s", err)
	}
	val := validator.New()
	engTs := eng.New()
	uni := ut.New(engTs, engTs)
	ts, _ := uni.GetTranslator("en")
	err = en.RegisterDefaultTranslations(val, ts)
	if err != nil {
		return nil, fmt.Errorf("failed to init service. Error: %s", err)
	}
	srv.handlers, err = handlers.NewHttpHandlers(rep, val, ts)
	if err != nil {
		return nil, fmt.Errorf("failed to init http handlers. Error: %s", err)
	}
	return &srv, nil
}

func (s *AuthService) Run() error {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/sign_in", s.handlers.SignIn)
		api.POST("/sign_up", s.handlers.SignUp)
	}

	err := router.Run(viper.GetString("server.host"))
	if err != nil {
		return fmt.Errorf("failed to start auth server. Error: %s", err)
	}
	return nil
}
