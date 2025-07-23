package services

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
)

type QuizService interface {
	GetAllQuizzes(ctx context.Context, filters web.FilterQuizRequest) ([]web.QuizResponse, error)
	StartQuiz(ctx context.Context, quizID uint, userID string) (web.QuizQuestionResponse, error)
	SubmitAnswer(ctx context.Context, request web.SubmitAnswerRequest) (web.SubmitAnswerResponse, error)
	//GetQuizAttemptResult(ctx context.Context, attemptID uint) (web.QuizAttemptResponse, error)
	SubmitDrawingAnswer(ctx context.Context, request web.SubmitDrawingRequest) (web.SubmitAnswerResponse, error)
	GetQuizzesByLessonID(ctx context.Context, lessonID uint, userID string) ([]web.QuizResponse, error)
}
