package main

import (
	"simpleTODO/internal/config"
	"simpleTODO/internal/service"
)

func main() {
	err := config.InitViper()
	if err != nil {
		panic(err)
	}
	srv, err := service.NewTodoServer()
	if err != nil {
		panic(err)
	}
	srv.Run()
}
