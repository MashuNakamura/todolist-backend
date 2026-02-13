package config

import (
	"fmt"
	"log"
	"os"

	"github.com/MashuNakamura/todolist-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	// Setup Environment Variable
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Connect to Database
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Connected to Database!")

	log.Println("Running Migrations...")

	// Migrate the every schema here to create the table
	err = DB.AutoMigrate(&models.User{}, &models.Task{}, &models.Category{})

	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Migrations success! Tables created.")
}
