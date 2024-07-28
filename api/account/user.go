package account

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"gsm/model"
	accountmodel "gsm/model/account"
	"gsm/pkg/cache"
	"gsm/pkg/errors"
	"gsm/pkg/util/hashext"
	"gsm/pkg/util/timeutil"
)

// userService defines the implementation of UserService interface
type userService struct {
	redisClient   cache.Client
	mongodbClient *mongo.Database
	jwtSecret     string
}

// UserService defines the user service interface
type UserService interface {
	CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error)
}

// NewUserService init the user service
func NewUserService(redisClient cache.Client, mongodbClient *mongo.Database, jwtSecret string) UserService {
	return &userService{
		redisClient:   redisClient,
		mongodbClient: mongodbClient,
		jwtSecret:     jwtSecret,
	}
}

func (impl *userService) CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	// check if email already exist
	duplicateEmail := accountmodel.Users{}
	cond := &model.MongoModelCond{
		MFilters: []bson.M{{"email": req.Email}},
	}
	if err := duplicateEmail.GetAll(ctx, impl.mongodbClient, cond); err != nil {
		return nil, errors.Errorf("failed to query users by cond %+v: %v", cond, err)
	} else if len(duplicateEmail) > 0 {
		return nil, errors.Errorf("email already registered")
	}

	// encrypted password and create user
	now := time.Now()
	encryptedPW, err := hashext.BcryptPassword(req.Password)
	if err != nil {
		return nil, errors.Errorf("failed to encrypted password: %v", err)
	}
	user := &accountmodel.User{
		ID:                ulid.Make().String(),
		Name:              req.Name,
		EncryptedPassword: encryptedPW,
		Email:             req.Email,
		CreatedAt:         timeutil.ConvertUTCTimeISOString(now),
		UpdatedAt:         timeutil.ConvertUTCTimeISOString(now),
	}
	if err := user.Create(ctx, impl.mongodbClient); err != nil {
		return nil, errors.Errorf("failed to create user: %v", err)
	}
	return &CreateUserResponse{}, nil
}
