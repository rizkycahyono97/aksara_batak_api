package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/rizkycahyono97/aksara_batak_api/config"
	"github.com/rizkycahyono97/aksara_batak_api/controllers"
	"github.com/rizkycahyono97/aksara_batak_api/repositories"
	"github.com/rizkycahyono97/aksara_batak_api/routes"
	"github.com/rizkycahyono97/aksara_batak_api/services"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//logger initialize
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	//load env
	err := godotenv.Load(".env")
	if err != nil {
		logger.Error("Error loading .env file")
	}

	//initialize other dependencies
	config.InitDB()
	validate := validator.New()

	//dependency injection
	authRepo := repositories.NewAuthRepository(config.DB)
	authService := services.NewAuthService(authRepo, validate, logger)
	authController := controllers.NewAuthController(authService, logger)

	//initialize fiber
	app := fiber.New()
	routes.SetupRoutes(app, authController)

	//gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// goroutine
	go func() {
		port := os.Getenv("APP_PORT")
		if port == "" {
			port = "8080"
		}
		logger.Info("Server is Running kudasai", "port", port)
		if err := app.Listen("0.0.0.0:" + port); err != nil {
			logger.Error("Failed to start server", "err", err)
			os.Exit(1)
		}
	}()
	//tunggu sinyal shutdown
	<-quit
	logger.Info("Shutting down server...")

	//memberi waktu 5 detik untuk menyelesaikan request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//shutdown
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to Shutdown", "error", err)
	}
	logger.Info("Server exiting")
}
