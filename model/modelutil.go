package model

import (
	"context"
	"gsm/pkg/errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Get(ctx context.Context, client *mongo.Database, id string, model MongoModel) error {
	if model == nil {
		return errors.NewErrorf(errors.InternalServerError, "model is empty")
	}

	collection := client.Collection(model.CollectionName())
	err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&model)
	if err != nil {
		return errors.NewErrorf(errors.NotFound, "failed to query by id %v: %v", id, err)
	}

	return nil
}

func Create(ctx context.Context, client *mongo.Database, model MongoModel) error {
	if model == nil {
		return errors.NewErrorf(errors.InternalServerError, "model is empty")
	}

	collection := client.Collection(model.CollectionName())
	_, err := collection.InsertOne(ctx, model)
	if err != nil {
		return errors.NewErrorf(errors.InternalServerError, "failed to insert model: %v", err)
	}

	return nil
}

func Update(ctx context.Context, client *mongo.Database, model MongoModel, filters *MongoModelCond, updates *MongoModelCond) (int64, error) {
	if model == nil {
		return 0, errors.NewErrorf(errors.InternalServerError, "model is empty")
	}

	collection := client.Collection(model.CollectionName())
	result, err := collection.UpdateMany(ctx, filters.ConstructQuery(), updates.ConstructQuery())
	if err != nil {
		return 0, errors.NewErrorf(errors.InternalServerError, "failed to update models: %v", err)
	}

	return result.ModifiedCount, nil
}

func Delete(ctx context.Context, client *mongo.Database, model MongoModel, filters *MongoModelCond) (int64, error) {
	if model == nil {
		return 0, errors.NewErrorf(errors.InternalServerError, "model is empty")
	}

	collection := client.Collection(model.CollectionName())
	result, err := collection.DeleteMany(ctx, filters.ConstructQuery())
	if err != nil {
		return 0, errors.NewErrorf(errors.InternalServerError, "failed to delete models: %v", err)
	}

	return result.DeletedCount, nil
}

func GetAll(ctx context.Context, client *mongo.Database, collectionName string, models MongoModels, filters *MongoModelCond) error {
	if models == nil {
		return errors.New("model is empty")
	}

	collection := client.Collection(collectionName)

	cursor, err := collection.Find(ctx, filters.ConstructQuery())
	if err != nil {
		return errors.NewErrorf(errors.InternalServerError, "failed to query models: %v", err)
	}
	defer cursor.Close(ctx)

	sliceValue := reflect.Indirect(reflect.ValueOf(models))
	elemType := sliceValue.Type().Elem()

	for cursor.Next(ctx) {
		elem := reflect.New(elemType).Interface()
		if err := cursor.Decode(elem); err != nil {
			return errors.NewErrorf(errors.InternalServerError, "failed to query models: %v", err)
		}
		sliceValue.Set(reflect.Append(sliceValue, reflect.ValueOf(elem).Elem()))
	}

	if err := cursor.Err(); err != nil {
		return errors.NewErrorf(errors.InternalServerError, "failed to handle query res: %v", err)
	}

	return nil
}

func CreateAll(ctx context.Context, client *mongo.Database, models MongoModels) error {
	return nil
}
