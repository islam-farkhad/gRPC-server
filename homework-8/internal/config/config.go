package config

import (
	"fmt"
	"homework-8/internal/pkg/db"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
)

func findEnvFile() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(currentDir, ".env")); err != nil {
			parentDir := filepath.Dir(currentDir)

			if parentDir == currentDir {
				break
			}
			currentDir = parentDir
			continue
		}

		return filepath.Join(currentDir, ".env"), nil
	}

	return "", fmt.Errorf(".env file not found")
}

func init() {
	envFilePath, err := findEnvFile()
	if err != nil {
		log.Fatal(".env file not found")
	}

	err = godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("Error loading .env file: %+v", err)
	}
}

// Config represents the configuration for Kafka.
type Config struct {
	dbConfig *db.Config
	apiPort  string
}

func newConfig() *Config {
	return &Config{
		dbConfig: getEnvDBConnectionConfig(),
		apiPort:  getAPIPort(),
	}
}

// GetDBConfig getter for dbConfig field
func (c Config) GetDBConfig() *db.Config {
	return c.dbConfig
}

// GetPort getter for apiPort field
func (c Config) GetPort() string {
	return c.apiPort
}

var (
	once     sync.Once
	instance *Config
)

// GetConfigs retrieves the Kafka configuration.
func GetConfigs() *Config {
	once.Do(func() {
		instance = newConfig()
	})
	return instance
}

func getAPIPort() string {
	apiPort := os.Getenv("HTTP_PORT")
	if apiPort == "" {
		log.Fatal("HTTP_PORT environment variable is not set")
	}
	return ":" + apiPort
}

func getEnvDBConnectionConfig() *db.Config {
	return &db.Config{
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("PASSWORD"),
		DBName:   os.Getenv("DBNAME"),
		SSLMode:  os.Getenv("SSLMODE"),
	}
}
