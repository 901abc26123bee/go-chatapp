package account

import (
	"context"
	"gsm/pkg/cache"
)

// authService defines the implementation of UserService interface
type authService struct {
	redisClient cache.Client
}

// AuthService defines the auth service interface
type AuthService interface {
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
}

// NewAuthService init the auth service
func NewAuthService(redisClient cache.Client) AuthService {
	return &authService{
		redisClient: redisClient,
	}
}

func (impl *authService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	return nil, nil
}
