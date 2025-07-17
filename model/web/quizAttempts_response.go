package web

import "time"

type QuizAttemptsResponse struct {
	//AttemptID   uint      `json:"attempt_id"`
	QuizID      uint      `json:"quiz_id"`
	QuizTitle   string    `json:"quiz_title"`
	Score       uint      `json:"score"`
	CompletedAt time.Time `json:"completed_at"`
}
