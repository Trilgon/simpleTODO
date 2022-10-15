package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"simpleTODO/internal/db"
	"simpleTODO/internal/db/postgresql"
)

type TodoServer struct {
	rep       db.Repository
	validator *validator.Validate
}

func NewTodoServer() (*TodoServer, error) {
	srv := TodoServer{}
	val := validator.New()
	rep, err := postgresql.NewTodoRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to init repository for TodoServer. Error: %s", err)
	}
	dbRep := db.Repository(rep)
	srv.rep = dbRep
	srv.validator = val
	return &srv, nil
}

func (s *TodoServer) Run() error {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/get_by_id", s.GetById)
		api.GET("/get_by_email", s.GetByEmail)
		api.GET("/search_by_text", s.SearchByText)
		api.GET("/mark_note", s.MarkNote)
		api.POST("/sign_up", s.SignUp)
		api.POST("sign_in", s.SignIn)
		api.POST("sign_out", s.SignOut)
		api.POST("/add_note", s.AddNote)
		api.POST("/delete_note", s.DeleteNotes)
		api.POST("/update_note", s.UpdateNote)
	}

	err := router.Run(viper.GetString("server.host"))
	if err != nil {
		return fmt.Errorf("failed to start server. Error: %s", err)
	}
	return nil
}
