package web

type TranslateDirection string

const (
	DirectionBatakAksaraToID         TranslateDirection = "batak_aksara_to_id"
	DirectionBatakAksaraToBatakLatin TranslateDirection = "batak_aksara_to_batak_latin"
	DirectionBatakLatinToBatakAksara TranslateDirection = "batak_latin_to_batak_aksara"
	DirectionBatakLatinToID          TranslateDirection = "batak_latin_to_id"
	DirectionIDToBatakAksara         TranslateDirection = "id_to_batak_aksara"
	DirectionIDToBatakLatin          TranslateDirection = "id_to_batak_latin"
)

type TranslateRequest struct {
	Text string `json:"text" validate:"required,min=1"`

	// Validasi memastikan field ini diisi dan nilainya harus salah satu dari dua pilihan.
	Direction TranslateDirection `json:"direction" validate:"required,oneof=batak_aksara_to_id batak_aksara_to_batak_latin batak_latin_to_batak_aksara batak_latin_to_id id_to_batak_aksara id_to_batak_latin"`
}

func (d TranslateDirection) IsValid() bool {
	switch d {
	case DirectionBatakAksaraToID, DirectionBatakAksaraToBatakLatin, DirectionBatakLatinToBatakAksara, DirectionBatakLatinToID, DirectionIDToBatakAksara, DirectionIDToBatakLatin:
		return true
	default:
		return false
	}
}
