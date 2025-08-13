package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig interface {
	Connect() *gorm.DB
}

type databaseConfig struct {
	db *gorm.DB
}

func NewDatabaseConfig(envConfig EnvConfig) DatabaseConfig {
	host := envConfig.DB_HOST
	port := envConfig.DB_PORT
	user := envConfig.DB_USER
	password := envConfig.DB_PASSWORD
	dbname := envConfig.DB_NAME
	sslmode := envConfig.SSLMODE
	timezone := envConfig.TIMEZONE

	// Build DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		host, port, user, password, dbname, sslmode, timezone,
	)

	// Connect
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return &databaseConfig{db: db}
}

func (c *databaseConfig) Connect() *gorm.DB {
	return c.db
}
