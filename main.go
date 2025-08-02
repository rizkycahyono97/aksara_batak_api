package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/rizkycahyono97/aksara_batak_api/config"
	"github.com/rizkycahyono97/aksara_batak_api/controllers"
	"github.com/rizkycahyono97/aksara_batak_api/middleware"
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

	// === DEPENDENCY INJECTION (URUTAN YANG BENAR) ===
	// 1. Inisialisasi semua REPOSITORY terlebih dahulu
	authRepo := repositories.NewUserRepository(config.DB)
	userProfileRepo := repositories.NewUserProfileRepository(config.DB)
	quizRepo := repositories.NewQuizRepository(config.DB)
	quizAttemptRepo := repositories.NewQuizAttemptRepository(config.DB)
	lessonRepo := repositories.NewLessonRepository(config.DB)
	chatRepo := repositories.NewChatHistoryRepository(config.DB)
	contactRepo := repositories.NewContactSubmissionsRepository(config.DB)

	// 2. Inisialisasi semua SERVICE
	authService := services.NewAuthService(authRepo, userProfileRepo, validate, logger)
	userProfileService := services.NewUserProfileService(authRepo, userProfileRepo, quizAttemptRepo, validate, logger)
	quizService := services.NewQuizService(quizRepo, quizAttemptRepo, validate, logger, userProfileRepo)
	leaderboardService := services.NewLeaderboardService(userProfileRepo, logger)
	lessonService := services.NewLessonService(lessonRepo, validate, logger)
	chatbotService := services.NewChatbotService(chatRepo, logger)
	translateService := services.NewTranslateService(logger)
	emailService := services.NewEmailService(logger)
	contactUsService := services.NewContactUsService(contactRepo, emailService, logger)

	// 3. Inisialisasi semua CONTROLLER
	authController := controllers.NewAuthController(authService, logger)
	userProfileController := controllers.NewUserProfileController(userProfileService, logger)
	quizController := controllers.NewQuizController(quizService, logger)
	leaderboardController := controllers.NewLeaderboardController(leaderboardService, logger)
	lessonsController := controllers.NewLessonController(lessonService, quizService, logger)
	chatbotController := controllers.NewChatbotController(chatbotService, logger, validate)
	translateController := controllers.NewTranslateController(translateService, logger)
	contactUsController := controllers.NewContactUsController(contactUsService, logger, validate)

	//initialize fiber,routes,static
	app := fiber.New()
	app.Static("/assets", "./public")
	app.Use(cors.New(middleware.SetupCors())) // CORS
	routes.SetupRoutes(
		app,
		authController,
		quizController,
		userProfileController,
		leaderboardController,
		lessonsController,
		chatbotController,
		translateController,
		contactUsController,
	)

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
