package services

import (
	"context"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
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
	Validate              *validator.Validate
	Log                   *slog.Logger
	QuizRepository        repositories.QuizRepository
	UserProfileRepository repositories.UserProfileRepository

	//untuk meyimpan sesi kuis yang aktif
	sessions     map[string]*QuizSession
	sessionMutex sync.RWMutex // mutex untuk melindungi akses ke map sessions
}

// constructor
func NewQuizService(
	quizRepository repositories.QuizRepository,
	validate *validator.Validate,
	log *slog.Logger,
	userProfileRepository repositories.UserProfileRepository,
) QuizService {
	return &QuizServiceImpl{
		QuizRepository:        quizRepository,
		UserProfileRepository: userProfileRepository,
		Validate:              validate,
		Log:                   log,
		sessions:              make(map[string]*QuizSession),
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

	//cek jika quiz memiliki pertanyaan
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

	//mengambil detail pertanyaan pertama + pilihan jawabanya
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

	//response akhir session struct + DTO response
	response := web.QuizQuestionResponse{
		SessionID:            sessionID,
		QuestionID:           question.ID,
		TotalQuestions:       len(questionIDs),
		CurrentQuestionIndex: 1,
		QuestionText:         question.QuestionText,
		Options:              options,
	}

	return response, nil
}

func (s *QuizServiceImpl) SubmitAnswer(ctx context.Context, request web.SubmitAnswerRequest) (web.SubmitAnswerResponse, error) {
	s.Log.InfoContext(ctx, "submit answer process started", "sessionID", request.SessionID)

	//validasi DTO request
	if err := s.Validate.Struct(request); err != nil {
		s.Log.ErrorContext(ctx, "validation error", "error", err, "request", request)
		return web.SubmitAnswerResponse{}, err
	}

	//ambil sesi kuis yang aktif dari memory
	s.sessionMutex.RLock()
	session, ok := s.sessions[request.SessionID]
	s.sessionMutex.RUnlock()
	if !ok {
		s.Log.InfoContext(ctx, "session does not exist", "sessionID", request.SessionID)
		return web.SubmitAnswerResponse{}, errors.New("invalid session ID")
	}

	//jika CurrentQuestionIndex sama dengan QuestionIDs maka quiz selesai
	if session.CurrentQuestionIndex >= len(session.QuestionIDs) {
		return web.SubmitAnswerResponse{}, errors.New("quiz has already been completed")
	}

	//memastikan user mengirimkan questionID sama dengan questionid di sesi
	currentQuestionID := session.QuestionIDs[session.CurrentQuestionIndex]
	if currentQuestionID != request.QuestionID {
		s.Log.WarnContext(ctx, "user answered the wrong question", "expected", currentQuestionID, "got", request.QuestionID)
		return web.SubmitAnswerResponse{}, errors.New("mismatched question ID")
	}

	//validasi jawaban + menambah jumlah score
	correctOptionID, err := s.QuizRepository.FindCorrectOptionID(ctx, currentQuestionID)
	if err != nil {
		s.Log.ErrorContext(ctx, "failed to find correct option ID", "optionID", request.OptionID)
		return web.SubmitAnswerResponse{}, err
	}
	isCorrect := request.OptionID == correctOptionID
	if isCorrect {
		s.Log.InfoContext(ctx, "option ID is match", "optionID", request.OptionID)
		session.CurrentScore += 10
	}
	session.CurrentQuestionIndex++ //increment index setiap selesai kirim jawaban

	//response dasar
	response := web.SubmitAnswerResponse{
		IsCorrect:       isCorrect,
		CorrectOptionID: correctOptionID,
	}

	//logika jika kuis jika sudah selesai
	if session.CurrentQuestionIndex >= len(session.QuestionIDs) {
		response.QuizFinished = true
		xpEarned := session.CurrentScore // 1 skor 1 xp
		response.FinalResult = &web.FinalResultResponse{
			FinalScore: xpEarned,
			XPEarned:   xpEarned,
		}

		//goroutine agar tidak memblokir response, dan kirim ke response
		go func() {
			//membuat context baru agar tidak terpengaruh oleh timout request
			bgCtx := context.Background()

			//assign struct kuis session ke quizAttempt
			attempt := &domain.QuizAttempts{UserID: session.UserID, QuizID: session.QuizID, Score: uint(session.CurrentScore)}
			if err := s.QuizRepository.CreateQuizAttempt(bgCtx, attempt); err != nil {
				s.Log.ErrorContext(bgCtx, "failed to create quiz attempt", "error", err)
			}

			//======================
			//menambah streak dan XP
			//======================
			// ambil userID
			profile, err := s.UserProfileRepository.FindUserProfileByID(ctx, attempt.UserID)
			if err != nil {
				s.Log.ErrorContext(bgCtx, "failed to find user profile for update", "error", err, "userID", session.UserID)
				return
			}

			//hitung streak
			var newStreak uint
			today := time.Now()
			yesterday := today.AddDate(0, 0, -1)

			if profile.LastActiveAt.Year() == yesterday.Year() && profile.LastActiveAt.YearDay() == yesterday.YearDay() { //jika kemarin aktif streak ditambah 1
				newStreak = profile.CurrentStreak + 1
			} else if profile.LastActiveAt.Year() == today.Year() || profile.LastActiveAt.YearDay() != today.YearDay() { //jika hari doang streak reset ke 1
				newStreak = 1
			} else {
				newStreak = profile.CurrentStreak // selain itu streak samakan dengan currentStreak
			}

			//repository untuk update streak dan xp
			if err := s.UserProfileRepository.UpdateXPAndStreak(bgCtx, session.UserID, xpEarned, newStreak, today); err != nil {
				s.Log.ErrorContext(bgCtx, "failed to update user xp and streak", "error", err, "userID", session.UserID)
			}
			s.Log.InfoContext(bgCtx, "user profile updated after quiz", "userID", session.UserID, "xp_earned", xpEarned, "new_streak", newStreak)
		}()

		//hapus session dari memory
		s.sessionMutex.Lock()
		delete(s.sessions, request.SessionID)
		s.sessionMutex.Unlock()
	} else {
		//untuk melanjutkan kuis selanjutnya, jika quiz_finished != true
		response.QuizFinished = false
		nextQuestionID := session.QuestionIDs[session.CurrentQuestionIndex]

		question, err := s.QuizRepository.FindQuestionWithOptions(ctx, nextQuestionID)
		if err != nil {
			return web.SubmitAnswerResponse{}, err
		}

		//DTO
		var options []web.QuestionOptionResponse
		for _, opt := range question.QuestionOptions {
			options = append(options, web.QuestionOptionResponse{
				ID:   opt.ID,
				Text: opt.OptionText,
			})
		}

		response.NextQuestion = &web.QuizQuestionResponse{
			SessionID:            request.SessionID,
			QuestionID:           question.ID,
			TotalQuestions:       len(session.QuestionIDs),
			CurrentQuestionIndex: session.CurrentQuestionIndex + 1,
			QuestionText:         question.QuestionText,
			Options:              options,
		}
	}

	return response, nil
}
