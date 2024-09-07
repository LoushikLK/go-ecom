package routes

import (
	routes_v1 "ecommerce/routes/v1"

	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) error {

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

	userRoute := v1.Group("/user") //api/v1/user
	routes_v1.InitAuthRoutes(userRoute.Group("/auth"))
	routes_v1.InitProfileRoutes(userRoute.Group("/profile"))
	routes_v1.InitAddressRoutes(userRoute.Group("/address"))

	return nil
}
