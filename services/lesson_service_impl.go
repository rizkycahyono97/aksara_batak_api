package services

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/repositories"
	"log/slog"
)

type LessonServiceImpl struct {
	LessonRepository repositories.LessonRepository
	Validate         *validator.Validate
	Log              *slog.Logger
}

func NewLessonService(
	lessonRepository repositories.LessonRepository,
	validate *validator.Validate,
	log *slog.Logger) LessonService {
	return &LessonServiceImpl{
		LessonRepository: lessonRepository,
		Validate:         validate,
		Log:              log,
	}
}

func (s *LessonServiceImpl) GetAllLessons(ctx context.Context) ([]web.LessonResponse, error) {
	s.Log.InfoContext(ctx, "get all lessons process started")

	//repository
	lessons, err := s.LessonRepository.FindAllLesson(ctx)
	if err != nil {
		s.Log.ErrorContext(ctx, "get all lessons failed", err)
		return nil, err
	}

	//DTO
	var lessonResponse []web.LessonResponse
	for _, lesson := range lessons {
		response := web.LessonResponse{
			ID:          int(lesson.ID),
			Title:       lesson.Title,
			Description: lesson.Description,
			IconURL:     lesson.IconURL,
		}
		lessonResponse = append(lessonResponse, response)
	}

	s.Log.InfoContext(ctx, "successfully retrieved lessons", "count", len(lessonResponse))
	return lessonResponse, nil
}
