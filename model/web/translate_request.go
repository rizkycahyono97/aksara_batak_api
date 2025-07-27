package web

type TranslateDirection string

const (
	DirectionBatakToID TranslateDirection = "BATAK_TO_ID"
	DirectionIDToBatak TranslateDirection = "ID_TO_BATAK"
)

type TranslateRequest struct {
	Text string `json:"text" validate:"required,min=1"`

	// Validasi memastikan field ini diisi dan nilainya harus salah satu dari dua pilihan.
	Direction TranslateDirection `json:"direction" validate:"required,oneof=BATAK_TO_ID ID_TO_BATAK"`
}

func (d TranslateDirection) IsValid() bool {
	switch d {
	case DirectionBatakToID, DirectionIDToBatak:
		return true
	default:
		return false
	}
}
