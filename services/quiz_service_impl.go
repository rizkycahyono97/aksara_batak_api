package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/repositories"
	"log/slog"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

// struct untuk meyimpan session quiz
type QuizSession struct {
	UserID               string
	QuizID               uint
	QuestionIDs          []uint
	CurrentQuestionIndex int
	CurrentScore         int
}

type QuizServiceImpl struct {
	Validate       *validator.Validate
	Log            *slog.Logger
	QuizRepository repositories.QuizRepository
	//UserProfileRepository repositories.UserProfileRepository

	//untuk meyimpan sesi kuis yang aktif
	sessions     map[string]*QuizSession
	sessionMutex sync.RWMutex // mutex untuk melindungi akses ke map sessions
}

// constructor
func NewQuizService(repo repositories.QuizRepository, validate *validator.Validate, log *slog.Logger) QuizService {
	return &QuizServiceImpl{
		QuizRepository: repo,
		Validate:       validate,
		Log:            log,
		sessions:       make(map[string]*QuizSession),
	}
}

// mengambil semua quiz yang ada
// menerima filter di untuk level, dialect, title
func (s *QuizServiceImpl) GetAllQuizzes(ctx context.Context, filters web.FilterQuizRequest) ([]web.QuizResponse, error) {
	s.Log.InfoContext(ctx, "get all quizzes process started", "dialect", filters.Dialect)

	//panggil repository
	quizzes, err := s.QuizRepository.FindAllQuizzes(ctx, filters)
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

// untuk memulai sebuah sesi kuis
func (s *QuizServiceImpl) StartQuiz(ctx context.Context, quizID uint, userID string) (web.QuizQuestionResponse, error) {
	s.Log.InfoContext(ctx, "start quiz", "quizID", quizID, "userID", userID)

	//ambil semua ID question untuk kuis ini
	questionIDs, err := s.QuizRepository.FindQuestionIDsByQuizID(ctx, quizID)
	if err != nil {
		s.Log.InfoContext(ctx, "failed to find questions IDs", "quizID", quizID, "userID", userID)
		return web.QuizQuestionResponse{}, err
	}

	//cek jika quiz ada
	if len(questionIDs) == 0 {
		s.Log.InfoContext(ctx, "failed to find questions IDs", "quizID", quizID, "userID", userID)
		return web.QuizQuestionResponse{}, errors.New("quiz not found or has no questions")
	}

	//mengacak urutan id question, supaya setiap questions bisa teracak
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(questionIDs), func(i, j int) {
		questionIDs[i], questionIDs[j] = questionIDs[j], questionIDs[i]
	})

	//membuat sesi kuis baru
	sessionID := uuid.NewString()
	session := &QuizSession{
		UserID:               userID,
		QuizID:               quizID,
		QuestionIDs:          questionIDs,
		CurrentQuestionIndex: 0,
		CurrentScore:         0,
	}

	//simpan sesi ke dalam map dengan mutex (perlindungan)
	s.sessionMutex.Lock()
	s.sessions[sessionID] = session
	s.sessionMutex.Unlock()
	s.Log.InfoContext(ctx, "new quiz session created", "sessionID", sessionID, "questionCount", len(questionIDs))

	//mengambil detail pertanyaan pertama
	firstQuestionID := questionIDs[0]
	question, err := s.QuizRepository.FindQuestionWithOptions(ctx, firstQuestionID)
	if err != nil {
		s.Log.InfoContext(ctx, "failed to find question details/answers", "error", err, "questionID", firstQuestionID)
		return web.QuizQuestionResponse{}, err
	}

	//DTO response
	var options []web.QuestionOptionResponse
	for _, opt := range question.QuestionOptions {
		options = append(options, web.QuestionOptionResponse{
			ID:   opt.ID,
			Text: opt.OptionText,
		})
	}

	response := web.QuizQuestionResponse{
		SessionID:            sessionID,
		TotalQuestions:       len(questionIDs),
		CurrentQuestionIndex: 1,
		QuestionText:         question.QuestionText,
		Options:              options,
	}

	return response, nil
}

func (s *QuizServiceImpl) SubmitAnswer(ctx context.Context, request web.SubmitAnswerRequest) (web.QuizQuestionResponse, error) {
	panic("error")
}
