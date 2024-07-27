package gormpsql

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"gsm/pkg/orm"
)

const (
	defaultMaxOpenConn = 50
	defaultMaxIdleConn = 10
	defaultLifetime    = time.Hour
)

// Initialize open a sql connection
func Initialize(configPath string) (orm.DB, error) {
	return InitializeWithEncryptedKey(configPath, "")
}

// InitializeWithCrpyoKey open a sql connection with the dbCryptoKey
func InitializeWithEncryptedKey(configPath, dbCryptoKey string) (orm.DB, error) {
	dbConfig, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get db config: %v", err.Error())
	}

	db, err := gorm.Open(postgres.Open(string(dbConfig)), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open db connection: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get db: %v", err)
	}
	sqlDB.SetMaxIdleConns(defaultMaxIdleConn)
	sqlDB.SetMaxOpenConns(defaultMaxOpenConn)
	sqlDB.SetConnMaxLifetime(defaultLifetime)

	return &adapter{db: db, key: dbCryptoKey}, err
}
