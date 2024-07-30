package router

import (
	"github.com/gofiber/fiber/v2"
)

var USER fiber.Router

func hello(c *fiber.Ctx) error {
	return c.SendString("Hello World!")
}

// SetupRoutes setups all the Routes
func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Healthy")
	})

	api := app.Group("/api")
	USER = api.Group("/user")

	SetupUserRoutes()

	api.Get("/", hello)
}
