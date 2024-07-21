package account

import (
	"context"
	"gsm/pkg/cache"
)

// userService defines the implementation of UserService interface
type userService struct {
	redisClient cache.Client
}

// UserService defines the user service interface
type UserService interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
}

// NewUserService init the user service
func NewUserService(redisClient cache.Client) UserService {
	return &userService{
		redisClient: redisClient,
	}
}

func (impl *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	return nil, nil
}
