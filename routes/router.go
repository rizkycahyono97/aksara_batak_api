package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizkycahyono97/aksara_batak_api/controllers"
)

func SetupRoutes(app *fiber.App, authController *controllers.AuthController) {
	//instance
	app = fiber.New()

	//intance middleware
	//jwtMiddleware := middleware.JWTMiddleware()

	//group
	publik := app.Group("/api/v1")
	publik.Post("/login", authController.Login)
	publik.Post("/register", authController.Register)

	//private route
	//private := app.Group("/api/v1", jwtMiddleware)
}
