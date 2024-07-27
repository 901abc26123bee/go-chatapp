package account

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"gsm/model"
)

type User struct {
	ID                string      `bson:"_id"`
	Name              string      `bson:"name,omitempty"`
	EncryptedPassword string      `bson:"encrypted_password,omitempty"`
	Email             string      `bson:"email,omitempty"`
	CreatedAt         interface{} `bson:"created_at,omitempty"`
	UpdatedAt         string      `bson:"updated_at,omitempty"`
}

func (m *User) CollectionName() string {
	return "users"
}

type Users []User

func (m *User) Get(ctx context.Context, client *mongo.Database) error {
	return model.Get(ctx, client, m.ID, m)
}

func (m *User) Create(ctx context.Context, client *mongo.Database) error {
	return model.Create(ctx, client, m)
}

func (m *User) Update(ctx context.Context, client *mongo.Database, filters *model.MongoModelCond, updates *model.MongoModelCond) (int64, error) {
	return model.Update(ctx, client, m, filters, updates)
}

func (m *User) Delete(ctx context.Context, client *mongo.Database, filters *model.MongoModelCond) (int64, error) {
	return model.Delete(ctx, client, m, filters)
}

func (ms *Users) GetAll(ctx context.Context, client *mongo.Database, filters *model.MongoModelCond) error {
	var user *User
	return model.GetAll(ctx, client, user.CollectionName(), ms, filters)
}

func (ms *Users) CreateAll(ctx context.Context, client *mongo.Database) error {
	return model.CreateAll(ctx, client, ms)
}
