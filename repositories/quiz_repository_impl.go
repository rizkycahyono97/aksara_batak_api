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
	return &QuizRepositoryImpl{db: db}
}

// method untuk mencari semua quiz, dan
// dengan beberapa filter seperti dialect, level, dan title
func (r *QuizRepositoryImpl) FindAllQuizzes(ctx context.Context, filters web.FilterQuizRequest) ([]domain.Quizzes, error) {
	var quizzes []domain.Quizzes

	//select * from
	query := r.db.WithContext(ctx).Model(domain.Quizzes{})

	//filters
	if filters.Dialect != "" {
		query = query.Where("dialect LIKE ?", "%"+filters.Dialect+"%")
	}
	if filters.Level != nil {
		query = query.Where("level LIKE ?", *filters.Level)
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

// method ini hanya mengambil satu id dari QuizID
// digunakan untuk memulai satu sesi kuis
func (r *QuizRepositoryImpl) FindQuestionIDsByQuizID(ctx context.Context, quizID uint) ([]uint, error) {
	//untuk hasil pertanyaan
	var questionsIDs []uint

	//query ke table questions, filter berdasarkan quiz_id, hanya id question bukan isinya
	// pluck -> SELECT id FROM questions WHERE quiz_id = {quizID};
	err := r.db.WithContext(ctx).
		Model(domain.Questions{}).
		Where("quiz_id = ?", quizID).
		Pluck("id", &questionsIDs).Error
	if err != nil {
		return nil, err
	}

	return questionsIDs, err
}

// method untuk mengambil detail lengkap dalam 1 pertanyaan
// termasuk semua jawaban dari 1 soal tersebut
func (r *QuizRepositoryImpl) FindQuestionWithOptions(ctx context.Context, questionID uint) (domain.Questions, error) {
	var question domain.Questions

	//query -> melakukan preload ke QuestionOptions untuk mengambil semua jawaban dalam 1 question
	// SELECT * FROM questions WHERE questionsID
	err := r.db.WithContext(ctx).
		Preload("QuestionOptions").
		Where("id", questionID).
		First(&question).Error
	if err != nil {
		return domain.Questions{}, err
	}

	return question, nil
}

// method untuk menemukan jawaban yang benar dari sebuah questions
func (r *QuizRepositoryImpl) FindCorrectOptionID(ctx context.Context, questionID uint) (uint, error) {
	var correctOptionID uint

	//query untuk menentukan is_correct == true
	err := r.db.WithContext(ctx).
		Model(&domain.QuestionOptions{}).
		Where("question_id = ? AND is_correct = ?", questionID, true).
		Pluck("id", &correctOptionID).Error
	if err != nil {
		return 0, err
	}

	return correctOptionID, nil
}

// menyimpan satu record question ke quiz_attempts
func (r *QuizRepositoryImpl) CreateQuizAttempt(ctx context.Context, attempt *domain.QuizAttempts) error {
	err := r.db.WithContext(ctx).Create(attempt).Error
	if err != nil {
		return err
	}

	return nil
}
