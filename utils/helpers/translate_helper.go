package helpers

import (
	"fmt"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"google.golang.org/genai"
)

// translate helper untuk prompt direction
func BuildTranslationPrompt(text string, direction web.TranslateDirection) string {
	switch direction {
	case web.DirectionBatakToID:
		return fmt.Sprintf("Terjemahkan teks Bahasa Batak Toba ini ke Bahasa Indonesia. Kembalikan HANYA hasil terjemahannya saja, tanpa penjelasan tambahan atau tanda kutip. Teks: \\\"%s\\\"\"", text)
	case web.DirectionIDToBatak:
		return fmt.Sprintf("Terjemahkan teks Bahasa Indonesia ini ke Bahasa Batak Toba. Kembalikan HANYA hasil terjemahannya saja, tanpa penjelasan tambahan atau tanda kutip. Teks: \"%s\"", text)
	default:
		return text
	}
}

// translate helper untuk parsing response
func ParseResponse(resp *genai.GenerateContentResponse) (string, error) {
	// 1. Cek jika response kosong
	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return "", fmt.Errorf("empty response from Gemini")
	}

	// 2. Iterasi semua parts
	for _, part := range resp.Candidates[0].Content.Parts {
		// 3. Akses field Text langsung dari *genai.Part
		if part != nil && part.Text != "" {
			return part.Text, nil
		}
	}

	// 4. Fallback jika tidak ada teks
	return "", fmt.Errorf("no text content found in response")
}
