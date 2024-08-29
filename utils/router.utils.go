package utils

import (
	"ecommerce/routes"

	"github.com/gofiber/fiber/v2"
)

func ApiRouteParser(app *fiber.App) error {

	app.Get("/", func(c *fiber.Ctx) error {
		c.Status(200)
		return c.JSON(fiber.Map{
			"message": "Welcome to the E-commerce API!",
			"success": true,
		})
	})

	api := app.Group("/api")
	v1 := api.Group("/v1", func(c *fiber.Ctx) error {
		c.Set("API-Version", "v1")
		return c.Next()
	})

	userRoute := v1.Group("/users") //api/v1/users
	routes.InitAuthRoutes(userRoute)

	return nil
}
