package account

import (
	"gsm/pkg/cache"
	"gsm/pkg/errors"
	"gsm/pkg/mdb"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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
	dbKey       string
}

// NewAccountController creates a new account controller
func NewAccountController(redisClient cache.Client, mongodb *mongo.Client, dbKey string, jwtSecret string) (AccountController, error) {
	return &accountController{
		userService: NewUserService(redisClient, mongodb.Database(mdb.DatabaseGSM), jwtSecret),
		authService: NewAuthService(redisClient, mongodb.Database(mdb.DatabaseGSM), jwtSecret),
		dbKey:       dbKey,
	}, nil
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

	// check request body
	if name := strings.TrimSpace(request.Name); name == "" {
		ctx.Error(errors.NewErrorf(errors.ParamIncorrect, "empty name"))
		return
	} else {
		request.Name = name
	}
	if email := strings.TrimSpace(request.Email); email == "" {
		// TODO: check email format
		ctx.Error(errors.NewErrorf(errors.ParamIncorrect, "empty email"))
		return
	} else {
		request.Email = email
	}
	if pw := strings.TrimSpace(request.Password); pw == "" {
		ctx.Error(errors.NewErrorf(errors.ParamIncorrect, "empty password"))
		return
	} else {
		request.Password = pw
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
	Email    string `json:"email"`
	Password string `json:"password"`
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

	// check request body
	if email := strings.TrimSpace(request.Email); email == "" {
		// TODO: check email format
		ctx.Error(errors.NewErrorf(errors.ParamIncorrect, "empty email"))
		return
	} else {
		request.Email = email
	}
	if pw := strings.TrimSpace(request.Password); pw == "" {
		ctx.Error(errors.NewErrorf(errors.ParamIncorrect, "empty password"))
		return
	} else {
		request.Password = pw
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
