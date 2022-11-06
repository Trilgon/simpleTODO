package handlers

import "github.com/gin-gonic/gin"

type AuthHandlers interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
}
