package mdb

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const DatabaseGSM string = "gsm"

// Initialize open a mongodb connection
func Initialize(ctx context.Context, configPath string) (*mongo.Client, error) {
	return InitializeWithEncryptedKey(ctx, configPath, "")
}

// InitializeWithCrpyoKey open a mongodb connection with the dbCryptoKey
func InitializeWithEncryptedKey(ctx context.Context, configPath, dbCryptoKey string) (*mongo.Client, error) {
	dbConfig, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get mondo db config: %v", err.Error())
	}
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().
		SetCompressors([]string{"zlib"}).
		ApplyURI(string(dbConfig)).
		SetServerAPIOptions(serverAPI).
		SetMaxPoolSize(16)
	// create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongo db: %v", err)
	}
	// send a ping to confirm a successful connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping to mongo db: %v", err)
	}

	return client, nil
}
