package database

import (
	"fmt"
	"log"

	"github.com/Ion-Stefan/saas-go-fiber/config"
	"github.com/Ion-Stefan/saas-go-fiber/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global variable that holds the connection to the database
var DB *gorm.DB

// ConnectDB initializes the database connection and runs migrations
func ConnectDB() error {
	var err error

	// PostgreSQL DSN format
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.Envs.DBHost,
		config.Envs.DBUser,
		config.Envs.DBPassword,
		config.Envs.DBName,
		config.Envs.DBPort,
	)

	// Open the connection to the database using PostgreSQL driver
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to the Database: %v\n", err)
		return err
	}

	// Run migrations
	log.Println("Running Migrations...")
	// Here add the models, from the models directory that you want to do migrations for
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		log.Printf("Failed to run migrations: %v\n", err)
		return err
	}

	log.Println("Connected Successfully to the Database")
	return nil
}
