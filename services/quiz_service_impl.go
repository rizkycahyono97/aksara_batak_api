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
	"math"
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
	QuizAttemptRepository repositories.QuizAttemptRepository
	UserProfileRepository repositories.UserProfileRepository

	//untuk meyimpan sesi kuis yang aktif
	sessions     map[string]*QuizSession
	sessionMutex sync.RWMutex // mutex untuk melindungi akses ke map sessions
}

// constructor
func NewQuizService(
	quizRepository repositories.QuizRepository,
	quizAttemptRepository repositories.QuizAttemptRepository,
	validate *validator.Validate,
	log *slog.Logger,
	userProfileRepository repositories.UserProfileRepository,
) QuizService {
	return &QuizServiceImpl{
		QuizRepository:        quizRepository,
		QuizAttemptRepository: quizAttemptRepository,
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
			ID:         opt.ID,
			OptionText: opt.OptionText,
			AksaraText: opt.AksaraText,
			ImageURL:   opt.ImageURL,
			AudioURL:   opt.AudioURL,
		})
	}

	//response akhir session struct + DTO response
	response := web.QuizQuestionResponse{
		SessionID:            sessionID,
		QuestionID:           question.ID,
		TotalQuestions:       len(questionIDs),
		CurrentQuestionIndex: 1,
		QuestionType:         question.QuestionType,
		QuestionText:         question.QuestionText,
		ImageURL:             question.ImageURL,
		AudioURL:             question.AudioURL,
		LottieURL:            question.LottieURL,
		Options:              options,
	}

	return response, nil
}

func (s *QuizServiceImpl) finishQuizSession(ctx context.Context, session *QuizSession) (*web.FinalResultResponse, error) {
	// 1. Ambil detail kuis untuk mendapatkan XpReward.
	quiz, err := s.QuizRepository.FindByID(ctx, session.QuizID)
	if err != nil {
		s.Log.ErrorContext(ctx, "failed to get quiz details for finalization", "error", err)
		return nil, err
	}

	// 2. Hitung skor kelulusan (70%).
	maxPossibleScore := float64(len(session.QuestionIDs) * 10) // Asumsi 10 poin per soal
	passingScore := uint(math.Ceil(maxPossibleScore * 0.7))

	// 3. Cek apakah pengguna lulus pada percobaan kali ini.
	userPassedThisAttempt := uint(session.CurrentScore) >= passingScore

	xpToUpdate := 0 // Default, tidak ada XP yang diberikan

	// 4. Hanya proses pemberian XP jika pengguna LULUS pada percobaan ini.
	if userPassedThisAttempt {
		// Periksa riwayat: Apakah ini kelulusan PERTAMA KALI?
		hasPassedBefore, err := s.QuizAttemptRepository.HasUserPassedQuizBefore(ctx, session.UserID, session.QuizID, passingScore)
		if err != nil {
			s.Log.ErrorContext(ctx, "failed to check past passing status", "error", err)
		} else if !hasPassedBefore {
			// Jika belum pernah lulus sebelumnya, berikan hadiah XP statis dari kuis.
			xpToUpdate = int(quiz.XpReward)
			s.Log.InfoContext(ctx, "First time passing quiz. Awarding static XP.", "userID", session.UserID, "quizID", session.QuizID, "xp_reward", xpToUpdate)
		} else {
			s.Log.InfoContext(ctx, "Already passed this quiz before. No XP awarded.", "userID", session.UserID, "quizID", session.QuizID)
		}
	}

	// 5. Buat respons final untuk dikirim ke pengguna.
	finalResult := &web.FinalResultResponse{
		FinalScore: session.CurrentScore,
		XPEarned:   xpToUpdate,
	}

	// 6. Jalankan proses penyimpanan ke DB di latar belakang.
	go func() {
		bgCtx := context.Background()

		// Selalu simpan riwayat pengerjaan.
		attempt := &domain.QuizAttempts{UserID: session.UserID, QuizID: session.QuizID, Score: uint(session.CurrentScore)}
		if err := s.QuizRepository.CreateQuizAttempt(bgCtx, attempt); err != nil {
			s.Log.ErrorContext(bgCtx, "failed to create quiz attempt", "error", err)
		}

		// Ambil profil untuk menghitung streak.
		profile, err := s.UserProfileRepository.FindUserProfileByID(bgCtx, session.UserID)
		if err != nil {
			s.Log.ErrorContext(bgCtx, "failed to find user profile for update", "error", err, "userID", session.UserID)
			return
		}

		// Hitung streak baru.
		var newStreak uint
		today := time.Now()
		yesterday := today.AddDate(0, 0, -1)
		if profile.LastActiveAt.Year() == yesterday.Year() && profile.LastActiveAt.YearDay() == yesterday.YearDay() {
			newStreak = profile.CurrentStreak + 1
		} else if profile.LastActiveAt.Year() != today.Year() || profile.LastActiveAt.YearDay() != today.YearDay() {
			newStreak = 1
		} else {
			newStreak = profile.CurrentStreak
		}

		// Panggil repository untuk update profil.
		if err := s.UserProfileRepository.UpdateXPAndStreak(bgCtx, session.UserID, xpToUpdate, newStreak, today); err != nil {
			s.Log.ErrorContext(bgCtx, "failed to update user xp and streak", "error", err, "userID", session.UserID)
		}
		s.Log.InfoContext(bgCtx, "user profile updated after quiz", "userID", session.UserID, "xp_earned", xpToUpdate, "new_streak", newStreak)
	}()

	return finalResult, nil
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

	//logika jika kuis jika sudah selesai & dan hanya menambahkan xp_earned ketika first ambil quiz
	if session.CurrentQuestionIndex >= len(session.QuestionIDs) {
		finalResult, err := s.finishQuizSession(ctx, session)
		if err != nil {
			s.Log.ErrorContext(ctx, "failed to finish quiz session", "error", err)
			return web.SubmitAnswerResponse{}, err
		}
		response.QuizFinished = true
		response.FinalResult = finalResult

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
				ID:         opt.ID,
				OptionText: opt.OptionText,
				AksaraText: opt.AksaraText,
				ImageURL:   opt.ImageURL,
				AudioURL:   opt.AudioURL,
			})
		}

		response.NextQuestion = &web.QuizQuestionResponse{
			SessionID:            request.SessionID,
			QuestionID:           question.ID,
			TotalQuestions:       len(session.QuestionIDs),
			CurrentQuestionIndex: session.CurrentQuestionIndex + 1,
			QuestionType:         question.QuestionType,
			QuestionText:         question.QuestionText,
			ImageURL:             question.ImageURL,
			AudioURL:             question.AudioURL,
			LottieURL:            question.LottieURL,
			Options:              options,
		}
	}

	return response, nil
}

// method khusus untuk submit drawing
func (s *QuizServiceImpl) SubmitDrawingAnswer(ctx context.Context, request web.SubmitDrawingRequest) (web.SubmitAnswerResponse, error) {
	s.Log.InfoContext(ctx, "submit drawing answer process started", "sessionID", request.SessionID)

	if err := s.Validate.Struct(request); err != nil {
		s.Log.ErrorContext(ctx, "validation error for drawing submission", "error", err)
		return web.SubmitAnswerResponse{}, err
	}

	s.sessionMutex.RLock()
	session, ok := s.sessions[request.SessionID]
	s.sessionMutex.RUnlock()
	if !ok {
		s.Log.WarnContext(ctx, "session does not exist for drawing submission", "sessionID", request.SessionID)
		return web.SubmitAnswerResponse{}, errors.New("invalid session ID")
	}

	if request.IsCorrect {
		session.CurrentScore += 10 // Skor standar
	}
	session.CurrentQuestionIndex++

	response := web.SubmitAnswerResponse{
		IsCorrect:       request.IsCorrect,
		CorrectOptionID: 0,
	}

	if session.CurrentQuestionIndex >= len(session.QuestionIDs) {
		// Panggil helper terpusat yang SAMA untuk menangani akhir kuis
		finalResult, err := s.finishQuizSession(ctx, session)
		if err != nil {
			return web.SubmitAnswerResponse{}, err
		}
		response.QuizFinished = true
		response.FinalResult = finalResult

		// Hapus sesi dari memori
		s.sessionMutex.Lock()
		delete(s.sessions, request.SessionID)
		s.sessionMutex.Unlock()
	} else {
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
				ID:         opt.ID,
				OptionText: opt.OptionText,
				AksaraText: opt.AksaraText,
				ImageURL:   opt.ImageURL,
				AudioURL:   opt.AudioURL,
			})
		}

		response.NextQuestion = &web.QuizQuestionResponse{
			SessionID:            request.SessionID,
			QuestionID:           question.ID,
			TotalQuestions:       len(session.QuestionIDs),
			CurrentQuestionIndex: session.CurrentQuestionIndex + 1,
			QuestionType:         question.QuestionType,
			QuestionText:         question.QuestionText,
			ImageURL:             question.ImageURL,
			AudioURL:             question.AudioURL,
			LottieURL:            question.LottieURL,
			Options:              options,
		}
	}

	return response, nil
}

func (s *QuizServiceImpl) GetQuizzesByLessonID(ctx context.Context, lessonID uint, userID string) ([]web.QuizResponse, error) {
	s.Log.InfoContext(ctx, "get quizzes by lesson ID process started", "lessonID", lessonID)

	//repository utuk mendapatkan semua kuis dalam pelajaran tertentu
	quizzesInLesson, err := s.QuizRepository.FindAllByLessonIDWithQuestionCount(ctx, lessonID)
	if err != nil {
		s.Log.ErrorContext(ctx, "failed to find quizzes with question count", "error", err, "lessonID", lessonID)
		return nil, err
	}

	//repository untuk dapat higherScore
	highestScoreMap, err := s.QuizAttemptRepository.FindHighestScoresByUserID(ctx, userID)
	if err != nil {
		s.Log.ErrorContext(ctx, "failed to find user highest scores", "error", err, "userID", userID)
		return nil, err
	}

	//DTO
	var quizResponse []web.QuizResponse
	for _, quiz := range quizzesInLesson {
		//hitung score maximum dan skor kelulusan (70%)
		maxPossibleScore := float64(quiz.QuestionCount * 10)
		passingScore := uint(math.Ceil(maxPossibleScore * 0.7))

		//ambil score tertinggi
		userHighestScore := highestScoreMap[quiz.ID]

		isCompleted := userHighestScore >= passingScore

		response := web.QuizResponse{
			ID:          quiz.ID,
			Title:       quiz.Title,
			Description: quiz.Description,
			Level:       strconv.Itoa(int(quiz.Level)),
			Dialect:     quiz.Dialect,
			IsCompleted: isCompleted,
		}
		quizResponse = append(quizResponse, response)
	}

	s.Log.InfoContext(ctx, "succesfully to get quizzes by lesson ID process finished", "lessonID", lessonID)
	return quizResponse, err
}
