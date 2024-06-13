package account_service

import (
	"context"
	"gsm/pkg/orm"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AccountController is the interface for account api
type AccountController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	GetUsers(ctx *gin.Context)
}

// accountController defines the implementation of AccountController interface
type accountController struct {
	db orm.DB
}

// NewAccountController creates a new account service
func NewAccountController(ctx context.Context, db orm.DB) (AccountController, error) {
	return &accountController{}, nil
}

func(impl *accountController) Register(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func(impl *accountController) Login(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func(impl *accountController) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func(impl *accountController) GetUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func(impl *accountController) GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}