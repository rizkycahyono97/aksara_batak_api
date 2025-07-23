package services

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
)

type LessonService interface {
	GetAllLessons(ctx context.Context) ([]web.LessonResponse, error)
}
