package services

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/repositories"
	"log/slog"
	"strconv"
)

type QuizServiceImpl struct {
	Repo     repositories.QuizRepository
	Validate *validator.Validate
	Log      *slog.Logger
}

func NewQuizService(repo repositories.QuizRepository, validate *validator.Validate, log *slog.Logger) QuizService {
	return &QuizServiceImpl{
		Repo:     repo,
		Validate: validate,
		Log:      log,
	}
}

// mengambil semua quiz yang ada
// menerima filter di untuk level, dialect, title
func (s QuizServiceImpl) GetAllQuizzes(ctx context.Context, filters web.FilterQuizRequest) ([]web.QuizResponse, error) {
	s.Log.InfoContext(ctx, "get all quizzes process started", "dialect", filters.Dialect)

	//panggil repository
	quizzes, err := s.Repo.FindAllQuizzes(ctx, filters)
	if err != nil {
		s.Log.ErrorContext(ctx, "failed to find quizzes", "err", err)
		return nil, err
	}

	//assign DTO
	var quizResponse []web.QuizResponse
	for _, quiz := range quizzes {
		quizResponse = append(quizResponse, web.QuizResponse{
			ID:          quiz.ID,
			Title:       quiz.Title,
			Description: quiz.Description,
			Level:       strconv.Itoa(int(quiz.Level)),
			Dialect:     quiz.Dialect,
		})
	}
	s.Log.InfoContext(ctx, "get all quizzes process finished", "response", quizResponse)

	return quizResponse, nil
}
