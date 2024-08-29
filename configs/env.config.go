package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DB_HOST          string
	DB_PORT          string
	DB_NAME          string
	DB_USER          string
	DB_PASS          string
	JWT_SECRET_KEY   string
	GO_ENV           string
	MEMECACHE_SERVER string
}

func AppEnv() EnvConfig {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}

	return EnvConfig{
		DB_HOST:          os.Getenv("POSTGRES_HOST"),
		DB_PORT:          os.Getenv("POSTGRES_PORT"),
		DB_NAME:          os.Getenv("POSTGRES_DB"),
		DB_USER:          os.Getenv("POSTGRES_USER"),
		DB_PASS:          os.Getenv("POSTGRES_PASSWORD"),
		JWT_SECRET_KEY:   os.Getenv("JWT_SECRET_KEY"),
		GO_ENV:           os.Getenv("GO_ENV"),
		MEMECACHE_SERVER: os.Getenv("MEMECACHE_SERVER"),
	}

}
