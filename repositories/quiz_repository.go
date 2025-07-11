package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
)

type QuizRepository interface {
	FindAllQuizzes(ctx context.Context, filters web.FilterQuizRequest) ([]domain.Quizzes, error)
	FindQuestionIDsByQuizID(ctx context.Context, quizID uint) ([]uint, error)
	FindQuestionWithOptions(ctx context.Context, questionID uint) (domain.Questions, error)
	FindCorrectOptionID(ctx context.Context, questionID uint) (uint, error)
	CreateQuizAttempt(ctx context.Context, attempt *domain.QuizAttempts) error
}
