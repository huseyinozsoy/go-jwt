package database

import (
	"fmt"
	"log"
	"os"

	"github.com/huseyinozsoy/go-jwt/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectToDB connects the server with database
func ConnectToDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file \n", err)
	}

	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("PSQL_USER"), os.Getenv("PSQL_PASS"), os.Getenv("PSQL_DBNAME"), os.Getenv("PSQL_PORT"))

	log.Print("Connecting to PostgreSQL DB...")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)

	}
	log.Println("Connected")

	// turned on the loger on info mode
	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Print("Running the migrations...")
	err = DB.AutoMigrate(&models.User{})

	if err != nil {
		log.Fatal("Failed to database migration. \n", err)
	}
}
