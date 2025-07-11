package services

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
)

type QuizService interface {
	GetAllQuizzes(ctx context.Context, filters web.FilterQuizRequest) ([]web.QuizResponse, error)
}
