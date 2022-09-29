package main

import (
	"simpleTODO/internal/config"
	"simpleTODO/internal/db/postgresql"
)

func main() {
	err := config.InitViper()
	if err != nil {
		panic(err)
	}
	rep, err := postgresql.NewTodoRepository()
	if err != nil {
		panic(err)
	}
	err = rep.SignUp("test@test.com", "enc_pas")
	if err != nil {
		panic(err)
	}
}
