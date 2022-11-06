package main

import (
	"todoServer/internal/config"
	"todoServer/internal/service"
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
	err = srv.Run()
	if err != nil {
		panic(err)
	}
}
