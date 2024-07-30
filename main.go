package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/huseyinozsoy/go-jwt/database"
	"github.com/huseyinozsoy/go-jwt/router"
	"github.com/joho/godotenv"
)

func LoadEnv(env string) {
	var envFile string
	switch env {
	case "production":
		envFile = ".env.production"
	case "staging":
		envFile = ".env.staging"
	case "test":
		envFile = ".env.test"
	default:
		envFile = ".env.test"
	}

	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		log.Fatalf("Error: %s file does not exist", envFile)
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Error loading %s file: %v", envFile, err)
	}
}

func CreateServer() *fiber.App {
	app := fiber.New()

	return app
}

func main() {
	app := CreateServer()

	database.ConnectToDB()

	app.Use(cors.New())

	// Use middlewares for each route
	app.Use(
		logger.New(), // add Logger middleware
	)

	env := os.Getenv("ENV")
	LoadEnv(env)

	log.Print(os.Getenv("ENAME"))

	router.SetupRoutes(app)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	log.Fatal(app.Listen(":8080"))
}
