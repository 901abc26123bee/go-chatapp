package account

import (
	"context"
	"gsm/pkg/cache"

	"go.mongodb.org/mongo-driver/mongo"
)

// authService defines the implementation of UserService interface
type authService struct {
	redisClient   cache.Client
	mongodbClient *mongo.Database
}

// AuthService defines the auth service interface
type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	Logout(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	RefreshToken(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}

// NewAuthService init the auth service
func NewAuthService(redisClient cache.Client, mongodbClient *mongo.Database) AuthService {
	return &authService{
		redisClient:   redisClient,
		mongodbClient: mongodbClient,
	}
}

func (impl *authService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	return nil, nil
}

func (impl *authService) Logout(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	return nil, nil
}

func (impl *authService) RefreshToken(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	return nil, nil
}
