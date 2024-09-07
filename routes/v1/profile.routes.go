package routes_v1

import (
	"ecommerce/controllers"
	"ecommerce/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitProfileRoutes(router fiber.Router) {
	router.Get("/", middlewares.IsAuthenticated, controllers.GetProfile)
	router.Put("/", middlewares.IsAuthenticated, controllers.UpdateProfile)
	router.Put("/update-email", middlewares.IsAuthenticated, controllers.UpdateEmail)
	router.Get("/verify-email", controllers.VerifyEmail) //This is get because user can verify by simply redirect to the browser
}
