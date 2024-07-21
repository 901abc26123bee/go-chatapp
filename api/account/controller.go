package account

import (
	"gsm/pkg/cache"
	"gsm/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AccountController is the interface for account api
type AccountController interface {
	CreateUser(ctx *gin.Context)
	Login(ctx *gin.Context)
	Logout(ctx *gin.Context)
	GetUser(ctx *gin.Context)
	GetUsers(ctx *gin.Context)
}

// accountController defines the implementation of AccountController interface
type accountController struct {
	userService UserService
	authService AuthService
}

// NewAccountController creates a new account controller
func NewAccountController(redisClient cache.Client) (AccountController, error) {
	return &accountController{
		userService: NewUserService(redisClient),
		authService: NewAuthService(redisClient),
	}, nil
}

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
}

type CreateUserResponse struct {
	// empty
}

func (impl *accountController) CreateUser(ctx *gin.Context) {
	// bind request body to struct
	request := &CreateUserRequest{}
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.Error(errors.NewErrorf(errors.ParamIncorrect,
			"failed to parse CreateUser request body: %v", err))
		return
	}

	// get response from service
	response, err := impl.userService.CreateUser(ctx, request)
	if err != nil {
		ctx.Error(err)
		return
	}

	// write response to json
	ctx.JSON(http.StatusOK, response)
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	Token string
}

func (impl *accountController) Login(ctx *gin.Context) {
	// bind request body to struct
	request := &LoginRequest{}
	if err := ctx.ShouldBind(&request); err != nil {
		ctx.Error(errors.NewErrorf(errors.ParamIncorrect,
			"failed to parse CreateUser request body: %v", err))
		return
	}

	// get response from service
	response, err := impl.authService.Login(ctx, request)
	if err != nil {
		ctx.Error(err)
		return
	}

	// write response to json
	ctx.JSON(http.StatusOK, response)
}

func (impl *accountController) Logout(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func (impl *accountController) GetUser(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}

func (impl *accountController) GetUsers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello, World!",
	})
}
