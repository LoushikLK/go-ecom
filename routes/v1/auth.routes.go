package routes_v1

import (
	"ecommerce/controllers"
	"ecommerce/middlewares"

	"github.com/gofiber/fiber/v2"
)

func InitAuthRoutes(router fiber.Router) {
	router.Post("/register", controllers.Register)
	router.Post("/resend-verification", controllers.ResendVerificationOtp)
	router.Post("/register-verify", controllers.VerifyAccountRegistration)
	router.Post("/login", controllers.Login)
	router.Post("/login-verify", controllers.LoginVerifyOTP)
	router.Get("/generate-token", controllers.GenerateToken)
	router.Put("/logout", middlewares.IsAuthenticated, controllers.LogoutUser)
}
