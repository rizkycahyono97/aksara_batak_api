package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
)

type LessonRepository interface {
	FindAllLesson(ctx context.Context) ([]domain.Lessons, error)
}
