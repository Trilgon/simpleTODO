package main

import (
	"authServer/internal/config"
	"authServer/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	err := config.InitViper()
	if err != nil {
		logrus.Panicf("failed to init configs through viper. Error: %s", err)
	}
	authService, err := service.NewAuthService()
	if err != nil {
		logrus.Panic(err)
	}
	err = authService.Run()
	if err != nil {
		logrus.Panic(err)
	}
}
