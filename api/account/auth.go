package account

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"gsm/model"
	accountmodel "gsm/model/account"
	"gsm/pkg/cache"
	"gsm/pkg/errors"
	"gsm/pkg/util/hashext"
	"gsm/pkg/util/jwtext"
)

// authService defines the implementation of UserService interface
type authService struct {
	redisClient   cache.Client
	mongodbClient *mongo.Database
	jwtSecret     string
}

// AuthService defines the auth service interface
type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	RefreshToken(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}

// NewAuthService init the auth service
func NewAuthService(redisClient cache.Client, mongodbClient *mongo.Database, jwtSecret string) AuthService {
	return &authService{
		redisClient:   redisClient,
		mongodbClient: mongodbClient,
		jwtSecret:     jwtSecret,
	}
}

func (impl *authService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// check user exited with matched password
	users := accountmodel.Users{}
	cond := &model.MongoModelCond{
		MFilters: []bson.M{{"email": req.Email}},
	}
	if err := users.GetAll(ctx, impl.mongodbClient, cond); err != nil {
		return nil, errors.Errorf("failed to query users by cond %+v: %v", cond, err)
	} else if len(users) == 0 {
		return nil, errors.NewErrorf(errors.NotFound, "failed to find user by email")
	}
	user := users[0]
	if err := hashext.ComparePasswords(user.EncryptedPassword, req.Password); err != nil {
		return nil, errors.Errorf("wrong password: %v", err)
	}

	// TODO: login failed count limit

	token, _, err := jwtext.CreateAccessToken(user.ID, impl.jwtSecret, nil, jwtext.AccessTokenDuration)
	if err != nil {
		return nil, errors.Errorf("failed to create access token: %v", err)
	}

	// set user online in redis
	key := fmt.Sprintf("online:%s", user.ID)
	if err = impl.redisClient.Set(ctx, key, true, 600*time.Second); err != nil {
		return nil, errors.Errorf("failed to set user online in redis: %v", err)
	}

	// TODO: create refresh token and store

	return &LoginResponse{Token: token}, nil
}

func (impl *authService) Logout(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// TODO: delete refresh token
	// TODO: delete all subscription for user
	return nil, nil
}

func (impl *authService) RefreshToken(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	return nil, nil
}
