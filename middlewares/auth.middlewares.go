package middlewares

import (
	"ecommerce/helpers"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Authorization header not found",
			"success": false,
		})
	}

	authHeader = strings.Replace(authHeader, "Bearer ", "", -1)
	validToken, err := helpers.ParseToken(authHeader)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": strings.Split(err.Error(), "token has invalid claims: ")[1],
			"success": false,
		})
	}

	if validToken.UserId == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
			"success": false,
		})
	}

	c.Locals("userId", validToken.UserId)
	c.Locals("email", validToken.Email)
	c.Locals("mobile", validToken.Mobile)

	return c.Next()

}
