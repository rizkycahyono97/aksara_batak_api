package web

type FilterQuizRequest struct {
	Dialect string `json:"dialect" validate:"omitempty;max=20"`
	Level   *uint  `json:"level" validate:"omitempty,min=0"`
	Title   string `json:"title" validate:"omitempty,max=100"`
}
