package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"gorm.io/gorm"
)

type QuizRepositoryImpl struct {
	db *gorm.DB
}

func NewQuizRepository(db *gorm.DB) QuizRepository {
	return &QuizRepositoryImpl{db}
}

func (r *QuizRepositoryImpl) FindAllQuizzes(ctx context.Context, filters web.FilterQuizRequest) ([]domain.Quizzes, error) {
	var quizzes []domain.Quizzes

	//select * from
	query := r.db.WithContext(ctx).Model(domain.Quizzes{})

	//filters
	if filters.Dialect == "" {
		query = query.Where("quizzes LIKE ?", "%"+filters.Dialect+"%")
	}
	if filters.Level == nil {
		query = query.Where("quizzes LIKE ?", *filters.Level)
	}
	if filters.Title != "" {
		query = query.Where("title LIKE ?", "%"+filters.Title+"%")
	}

	err := query.Find(&quizzes).Error
	if err != nil {
		return nil, err
	}

	return quizzes, nil
}
