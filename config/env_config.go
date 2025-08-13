package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DB_HOST                     string
	DB_PORT                     string
	DB_USER                     string
	DB_PASSWORD                 string
	DB_NAME                     string
	SSLMODE                     string
	TIMEZONE                    string
	SECRETKEY                   string
	ACCESS_TOKEN_EXPIRE_MINUTES int
}

func NewEnvConfig() EnvConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	expireMinutes, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRE_MINUTES"))
	if err != nil {
		log.Fatalf("Invalid ACCESS_TOKEN_EXPIRE_MINUTES value: %v", err)
	}

	return EnvConfig{
		DB_HOST:                     os.Getenv("DB_HOST"),
		DB_PORT:                     os.Getenv("DB_PORT"),
		DB_USER:                     os.Getenv("DB_USER"),
		DB_PASSWORD:                 os.Getenv("DB_PASSWORD"),
		DB_NAME:                     os.Getenv("DB_NAME"),
		SSLMODE:                     os.Getenv("SSLMODE"),
		TIMEZONE:                    os.Getenv("TIMEZONE"),
		SECRETKEY:                   os.Getenv("SECRETKEY"),
		ACCESS_TOKEN_EXPIRE_MINUTES: expireMinutes,
	}
}
