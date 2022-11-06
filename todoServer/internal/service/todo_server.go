package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/spf13/viper"
	"todoServer/internal/db"
	"todoServer/internal/db/postgresql"
)

type TodoServer struct {
	rep       db.TodoRepository
	validator *validator.Validate
	ts        ut.Translator
}

func NewTodoServer() (*TodoServer, error) {
	srv := TodoServer{}
	val := validator.New()
	rep, err := postgresql.NewTodoRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to init repository for TodoServer. Error: %s", err)
	}
	srv.rep = rep
	srv.validator = val
	eng := en.New()
	uni := ut.New(eng, eng)
	srv.ts, _ = uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(srv.validator, srv.ts)
	return &srv, nil
}

func (s *TodoServer) Run() error {
	router := gin.Default()

	api := router.Group("/api")
	api.Use(s.AuthMiddleware)
	{
		api.GET("/get_by_id", s.GetById)
		api.GET("/get_by_email", s.GetByEmail)
		api.POST("/search_by_text", s.SearchByText)
		api.PATCH("/mark_note", s.MarkNote)
		//api.POST("/sign_up", s.SignUp)
		//api.POST("sign_in", s.SignIn)
		//api.POST("sign_out", s.SignOut)
		api.POST("/add_note", s.AddNote)
		api.DELETE("/delete_notes", s.DeleteNotes)
		api.PATCH("/update_note", s.UpdateNote)
	}

	err := router.Run(viper.GetString("server.host"))
	if err != nil {
		return fmt.Errorf("failed to start server. Error: %s", err)
	}
	return nil
}
