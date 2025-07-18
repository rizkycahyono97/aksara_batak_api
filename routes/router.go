package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rizkycahyono97/aksara_batak_api/controllers"
	"github.com/rizkycahyono97/aksara_batak_api/middleware"
)

func SetupRoutes(
	app *fiber.App,
	authController *controllers.AuthController,
	quizController *controllers.QuizController,
	userProfileController *controllers.UserProfileController,
	leaderboardController *controllers.LeaderboardController,
) {
	//intance middleware
	jwtMiddleware := middleware.JWTMiddleware()

	//=============
	//public route
	//=============
	publik := app.Group("/api/v1")
	publik.Post("/login", authController.Login)
	publik.Post("/register", authController.Register)

	//leaderboard
	publik.Get("/leaderboard", leaderboardController.GetLeaderboards)

	//=============
	//private route
	//=============
	private := app.Group("/api/v1", jwtMiddleware)

	// quiz
	private.Get("/quizzes", quizController.GetAllQuizzes)
	private.Get("/quizzes/:quizID/start", quizController.StartQuiz)
	private.Post("/quizzes/submit", quizController.SubmitQuiz)
	private.Post("/quizzes/submit-drawing", quizController.SubmitDrawingAnswer)

	// userProfile
	private.Get("/users/profile", userProfileController.GetMyProfile)
	private.Put("/users/profile", userProfileController.UpdateMyProfile)
	private.Get("/users/profile/attempts", userProfileController.GetMyAttempts)
}
