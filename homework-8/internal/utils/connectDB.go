package utils

import (
	"context"
	"fmt"
	"homework-8/internal/config"
	"homework-8/internal/pkg/db"
	"log"
)

// ConnectDB connects to db and returns Database instance
func ConnectDB(ctx context.Context) *db.Database {
	database, err := db.NewDB(ctx, MakeDBConnStr(config.GetConfigs().GetDBConfig()))
	if err != nil {
		log.Fatal(err)
	}
	return database
}

// MakeDBConnStr helper function
func MakeDBConnStr(config *db.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
}
