package routes_v1

import (
	"ecommerce/controllers"
	"ecommerce/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitAddressRoutes(router fiber.Router) {
	router.Post("/", middlewares.IsAuthenticated, controllers.AddAddress)
	router.Put("/:addressId", middlewares.IsAuthenticated, controllers.UpdateAddress)
	router.Delete("/:addressId", middlewares.IsAuthenticated, controllers.DeleteAddress)
	router.Get("/", middlewares.IsAuthenticated, controllers.GetAllAddresses)
}
