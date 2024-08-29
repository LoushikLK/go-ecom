package routes

import (
	"ecommerce/controllers"
	"ecommerce/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitAuthRoutes(router fiber.Router) {
	router.Post("/register", controllers.Register)
	router.Post("/login", controllers.Login)
	router.Get("/generate-token", controllers.GenerateToken)
	router.Put("/logout", middlewares.IsAuthenticated, controllers.LogoutUser)
	router.Get("/profile", middlewares.IsAuthenticated, controllers.GetProfile)
	router.Put("/profile", middlewares.IsAuthenticated, controllers.UpdateProfile)
}
