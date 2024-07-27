package model

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// MongoModel defines an interface for general CRUD methods of mdb model
type MongoModel interface {
	CollectionName() string

	Get(ctx context.Context, client *mongo.Database) error

	Create(ctx context.Context, client *mongo.Database) error

	Update(ctx context.Context, client *mongo.Database, filters *MongoModelCond, updates *MongoModelCond) (int64, error)

	Delete(ctx context.Context, client *mongo.Database, filters *MongoModelCond) (int64, error)
}

// MongoModels defines an interface for general batch CRUD methods of mdb models
type MongoModels interface {
	GetAll(ctx context.Context, client *mongo.Database, filters *MongoModelCond) error

	CreateAll(ctx context.Context, client *mongo.Database) error
}

type MongoModelCond struct {
	DFilters []bson.D
	MFilters []bson.M
	EFilters []bson.E
	AFilters []bson.A
}

func (q *MongoModelCond) ConstructQuery() bson.D {
	result := bson.D{}

	// add bson.D filters
	for _, d := range q.DFilters {
		result = append(result, d...)
	}

	// add bson.M filters (convert to bson.D)
	for _, m := range q.MFilters {
		for k, v := range m {
			result = append(result, bson.E{Key: k, Value: v})
		}
	}

	// add bson.E filters
	for _, e := range q.EFilters {
		result = append(result, e)
	}

	// add bson.A filters (this part assumes you know where to place arrays)
	for i, a := range q.AFilters {
		key := fmt.Sprintf("arrayFilter%d", i)
		result = append(result, bson.E{Key: key, Value: a})
	}

	return result
}
