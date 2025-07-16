package web

type QuizResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Dialect     string `json:"dialect"`
}

// DTO untuk satu pilihan jawaban
// tidak ada is_correct untuk keamanan
type QuestionOptionResponse struct {
	ID         uint   `json:"id"`
	OptionText string `json:"option_text"`
	AksaraText string `json:"aksara_text"`
	ImageURL   string `json:"image_url"`
	AudioURL   string `json:"audio_url"`
}

// DTO yang dikirim saat kuis dimulai
// atau saat pertanyaan berikutnya ditampilkan
type QuizQuestionResponse struct {
	SessionID            string                   `json:"session_id"`
	QuestionID           uint                     `json:"question_id"`
	TotalQuestions       int                      `json:"total_questions"`
	CurrentQuestionIndex int                      `json:"current_question_index"`
	QuestionType         string                   `json:"question_type"`
	QuestionText         string                   `json:"question_text"`
	ImageURL             string                   `json:"image_url"`
	AudioURL             string                   `json:"audio_url"`
	LottieURL            string                   `json:"lottie_url"`
	Options              []QuestionOptionResponse `json:"options"`
}

// DTO untuk response answer setelah menjawan satu soal
type SubmitAnswerResponse struct {
	IsCorrect       bool                  `json:"is_correct"`
	CorrectOptionID uint                  `json:"correct_option_id"`
	QuizFinished    bool                  `json:"quiz_finished"`
	NextQuestion    *QuizQuestionResponse `json:"next_question,omitempty"`
	FinalResult     *FinalResultResponse  `json:"final_result,omitempty"`
}

type FinalResultResponse struct {
	FinalScore int `json:"final_score"`
	XPEarned   int `json:"xp_earned"`
}
