package routes

import (
	"tsilodot/controller"

	"github.com/gofiber/fiber/v3"
)

func SetupAuthRoutes(app fiber.Router, authController *controller.AuthController) {
	authRoute := app.Group("/auth")

	authRoute.Post("/register", authController.Register)
	authRoute.Post("/login", authController.Login)
}
