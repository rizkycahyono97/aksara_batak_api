package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/rizkycahyono97/aksara_batak_api/config"
)

// MIDDLEWARE JWT menggunakan jwtware
// docs -> https://docs.gofiber.io/contrib/jwt/#install
func JWTMiddleware() fiber.Handler {
	jwtSecret := []byte(config.GetEnv("JWT_SECRET_KEY", ""))

	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: "HS256",
			Key:    jwtSecret,
		},
		//successhandler akan menyimpan di c.Locals (context)
		SuccessHandler: func(c *fiber.Ctx) error {
			return c.Next()
		},
		//error handler jika tika valid
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err.Error() == "Missing or malformed JWT" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"message": "Missing or malformed JWT",
				})
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid or expired JWT",
			})
		},
	})
}
