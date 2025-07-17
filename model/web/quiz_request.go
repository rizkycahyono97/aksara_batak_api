package web

type FilterQuizRequest struct {
	Dialect string `json:"dialect" validate:"omitempty;max=20"`
	Level   *uint  `json:"level" validate:"omitempty,min=0"`
	Title   string `json:"title" validate:"omitempty,max=100"`
}

// DTO untuk request submit
type SubmitAnswerRequest struct {
	SessionID  string `json:"session_id" validate:"required"`
	QuestionID uint   `json:"question_id" validate:"required"`
	OptionID   uint   `json:"option_id" validate:"required"`
}

// SubmitDrawingRequest DTO untuk data yang dikirim user
// ketika question_type == drawing
type SubmitDrawingRequest struct {
	SessionID  string `json:"session_id" validate:"required"`
	QuestionID uint   `json:"question_id" validate:"required"`
	Score      int    `json:"score" validate:"gte=0,lte=100"`
}
