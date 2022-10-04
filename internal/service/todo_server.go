package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"simpleTODO/internal/db"
	"simpleTODO/internal/db/postgresql"
)

type TodoServer struct {
	rep *db.Repository
}

func NewTodoServer() (*TodoServer, error) {
	srv := TodoServer{}
	rep, err := postgresql.NewTodoRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to init repository for TodoServer. Error: %s", err)
	}
	dbRep := db.Repository(rep)
	srv.rep = &dbRep
	return &srv, nil
}

func (s *TodoServer) Run() {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/get_by_id")
		api.GET("/get_by_email")
		api.GET("/search_by_text")
		api.POST("/sign_up")
		api.POST("sign_in")
		api.POST("sign_out")
		api.POST("/add_note")
		api.POST("/delete_note")
		api.POST("/update_note")
		api.POST("/mark_note")
	}
}
