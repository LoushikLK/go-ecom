package routes_v1

import (
	"ecommerce/controllers"
	"ecommerce/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitAddressRoutes(router fiber.Router) {
	router.Post("/", middlewares.IsAuthenticated, controllers.GetProfile)
	router.Put("/:addressId", middlewares.IsAuthenticated, controllers.UpdateProfile)
	router.Delete("/", middlewares.IsAuthenticated, controllers.UpdateEmail)
	router.Get("/", controllers.VerifyEmail)
}
