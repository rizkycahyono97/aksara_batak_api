package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"gorm.io/gorm"
)

type LessonRepositoryImpl struct {
	DB *gorm.DB
}

func NewLessonRepository(db *gorm.DB) LessonRepository {
	return &LessonRepositoryImpl{
		DB: db,
	}
}

func (r LessonRepositoryImpl) FindAllLesson(ctx context.Context) ([]domain.Lessons, error) {
	var lessons []domain.Lessons

	//Order: Mengurutkan pelajaran berdasarkan 'order_index' secara menaik (ASC).
	//Find = SELECT * FROM lessons WHERE id IN (1,2,3,4)
	err := r.DB.WithContext(ctx).
		Order("order_index ASC").
		Find(&lessons).Error
	if err != nil {
		return nil, err
	}

	return lessons, err
}
