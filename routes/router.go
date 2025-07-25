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
	lessonsController *controllers.LessonController,
	chatbotController *controllers.ChatbotController,
) {
	//intance middleware
	jwtMiddleware := middleware.JWTMiddleware()

	//=============
	//public route
	//=============
	public := app.Group("/api/v1")
	public.Post("/login", authController.Login)
	public.Post("/register", authController.Register)

	//leaderboard
	public.Get("/leaderboard", leaderboardController.GetLeaderboards)

	//chatbot
	public.Post("/chatpub", chatbotController.HandlePublicChat)

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

	// lessons
	private.Get("/lessons", lessonsController.GetAllLessons)
	private.Get("/lessons/:lessonID/quizzes", lessonsController.GetQuizzesByLessonID)
}
