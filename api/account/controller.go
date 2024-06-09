package account_service

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AccountController is the interface for account api
type AccountController interface {
	CreateUser(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	LoginWithEmailPassword(ctx *gin.Context)
}

// accountController defines the implementation of AccountController interface
type accountController struct {
}

// NewAccountController creates a new account service
func NewAccountController(ctx context.Context) (AccountController, error) {
	return &accountController{}, nil
}

func(impl *accountController) CreateUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func(impl *accountController) LoginWithEmailPassword(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func(impl *accountController) GetUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}