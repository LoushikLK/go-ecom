package middlewares

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func ErrorMiddleware(app *fiber.App) error {

	app.Use(func(c *fiber.Ctx) error {
		// Call the next handler in the stack
		err := c.Next()
		if err != nil {
			// Log the error
			log.Println("Error:", err)

			// Send a consistent error response
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"message": "Internal Server Error",
			})
		}
		return nil
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Not Found",
			"success": false,
		})
	})

	return nil
}
