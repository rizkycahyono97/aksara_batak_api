package helpers

import (
	"fmt"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
)

// translate helper untuk prompt direction
func BuildTranslationPrompt(text string, direction web.TranslateDirection) string {
	var _ = map[string]map[rune]string{
		"h":  {'a': "ᯂ", 'i': "ᯂ᯦", 'u': "ᯂᯧ", 'e': "ᯂᯨ", 'o': "ᯂᯩ"},
		"b":  {'a': "ᯅ", 'i': "ᯅ᯦", 'u': "ᯅᯧ", 'e': "ᯅᯨ", 'o': "ᯅᯩ"},
		"k":  {'a': "ᯆ", 'i': "ᯆ᯦", 'u': "ᯆᯧ", 'e': "ᯆᯨ", 'o': "ᯆᯩ"},
		"d":  {'a': "ᯇ", 'i': "ᯇ᯦", 'u': "ᯇᯧ", 'e': "ᯇᯨ", 'o': "ᯇᯩ"},
		"g":  {'a': "ᯈ", 'i': "ᯈ᯦", 'u': "ᯈᯧ", 'e': "ᯈᯨ", 'o': "ᯈᯩ"},
		"j":  {'a': "ᯉ", 'i': "ᯉ᯦", 'u': "ᯉᯧ", 'e': "ᯉᯨ", 'o': "ᯉᯩ"},
		"l":  {'a': "ᯊ", 'i': "ᯊ᯦", 'u': "ᯊᯧ", 'e': "ᯊᯨ", 'o': "ᯊᯩ"},
		"m":  {'a': "ᯋ", 'i': "ᯋ᯦", 'u': "ᯋᯧ", 'e': "ᯋᯨ", 'o': "ᯋᯩ"},
		"n":  {'a': "ᯌ", 'i': "ᯌ᯦", 'u': "ᯌᯧ", 'e': "ᯌᯨ", 'o': "ᯌᯩ"},
		"ng": {'a': "ᯍ", 'i': "ᯍ᯦", 'u': "ᯍᯧ", 'e': "ᯍᯨ", 'o': "ᯍᯩ"},
		"p":  {'a': "ᯎ", 'i': "ᯎ᯦", 'u': "ᯎᯧ", 'e': "ᯎᯨ", 'o': "ᯎᯩ"},
		"r":  {'a': "ᯏ", 'i': "ᯏ᯦", 'u': "ᯏᯧ", 'e': "ᯏᯨ", 'o': "ᯏᯩ"},
		"s":  {'a': "ᯐ", 'i': "ᯐ᯦", 'u': "ᯐᯧ", 'e': "ᯐᯨ", 'o': "ᯐᯩ"},
		"t":  {'a': "ᯑ", 'i': "ᯑ᯦", 'u': "ᯑᯧ", 'e': "ᯑᯨ", 'o': "ᯑᯩ"},
		"w":  {'a': "ᯒ", 'i': "ᯒ᯦", 'u': "ᯒᯧ", 'e': "ᯒᯨ", 'o': "ᯒᯩ"},
		"y":  {'a': "ᯓ", 'i': "ᯓ᯦", 'u': "ᯓᯧ", 'e': "ᯓᯨ", 'o': "ᯓᯩ"},
	}

	switch direction {
	case web.DirectionBatakAksaraToID:
		return fmt.Sprintf(`Terjemahkan teks Aksara Batak Toba ini ke Bahasa Indonesia. Kembalikan HANYA hasil terjemahannya saja, tanpa penjelasan tambahan atau tanda kutip. Teks: %s`, text)
	case web.DirectionBatakAksaraToBatakLatin:
		return fmt.Sprintf("Terjemahkan teks Aksara Batak Toba ini ke Bahasa Batak Toba Latin. Kembalikan HANYA hasil terjemahannya saja, tanpa penjelasan tambahan atau tanda kutip. Teks: %s", text)
	case web.DirectionBatakLatinToID:
		return fmt.Sprintf("Terjemahkan teks Bahasa Batak Toba ini ke Bahasa Indonesia. Kembalikan HANYA hasil terjemahannya saja, tanpa penjelasan tambahan atau tanda kutip. Teks: %s", text)
	case web.DirectionBatakLatinToBatakAksara:
		return fmt.Sprintf("Terjemahkan teks Bahasa Batak Toba Latin ini ke Aksara Batak Toba. Kembalikan HANYA hasil terjemahannya saja, tanpa penjelasan tambahan atau tanda kutip. Teks: %s", text)
	case web.DirectionIDToBatakLatin:
		return fmt.Sprintf("Terjemahkan teks Bahasa Indonesia ini ke Bahasa Batak Toba. Kembalikan HANYA hasil terjemahannya saja, tanpa penjelasan tambahan atau tanda kutip. Teks: %s", text)
	case web.DirectionIDToBatakAksara:
		return fmt.Sprintf("Terjemahkan teks Bahasa Indonesia ini ke Aksara Batak Toba. Kembalikan HANYA hasil terjemahannya saja, tanpa penjelasan tambahan atau tanda kutip. Teks: %s", text)
	default:
		return text
	}
}

// translate helper untuk parsing response
//func ParseResponse(resp *genai.GenerateContentResponse) (string, error) {
//	// 1. Cek jika response kosong
//	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
//		return "", fmt.Errorf("empty response from Gemini")
//	}
//
//	// 2. Iterasi semua parts
//	for _, part := range resp.Candidates[0].Content.Parts {
//		// 3. Akses field Text langsung dari *genai.Part
//		if part != nil && part.Text != "" {
//			return part.Text, nil
//		}
//	}
//
//	// 4. Fallback jika tidak ada teks
//	return "", fmt.Errorf("no text content found in response")
//}
