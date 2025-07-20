package middleware

import (
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rizkycahyono97/aksara_batak_api/utils/helpers"
	"os"
)

func SetupCors() cors.Config {
	return cors.Config{
		AllowOrigins:     os.Getenv("CORS_ALLOWED_ORIGINS"),
		AllowMethods:     os.Getenv("CORS_ALLOWED_METHODS"),
		AllowHeaders:     os.Getenv("CORS_ALLOWED_HEADERS"),
		AllowCredentials: helpers.GetEnvBool("CORS_ALLOW_CREDENTIALS", false),
		MaxAge:           12 * 60 * 60, // 12 hours
	}
}
