package web

type TranslateResponse struct {
	OriginalText   string             `json:"original_text"`
	TranslatedText string             `json:"translated_text"`
	Direction      TranslateDirection `json:"direction"`
	Timestamp      string             `json:"timestamp"` // Format: RFC3339
}

type TranslateErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
