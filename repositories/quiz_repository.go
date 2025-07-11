package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
)

type QuizRepository interface {
	FindAllQuizzes(ctx context.Context, filters web.FilterQuizRequest) ([]domain.Quizzes, error)
}
