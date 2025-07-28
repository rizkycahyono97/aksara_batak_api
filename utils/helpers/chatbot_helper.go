package helpers

import (
	"fmt"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"google.golang.org/genai"
	"strings"
)

func BuildPrompt(newMsg string, histories []domain.ChatHistories) string {
	var sb strings.Builder

	sb.WriteString("Anda adalah asisten ahli budaya Batak. Berikut percakapan terakhir:\\n")

	//tambahkan histories
	for _, h := range histories {
		sb.WriteString(fmt.Sprintf("User: %s\nAnda: %s\n", h.Message, h.Reply))
	}

	// Tambahkan pesan baru
	sb.WriteString(fmt.Sprintf("\nUser: %s\nAnda: ", newMsg))

	return sb.String()
}

func GetGenerationConfig() *genai.GenerationConfig {
	return &genai.GenerationConfig{
		MaxOutputTokens: 500,
		StopSequences:   []string{"User:"}, // Berhenti jika ada input user baru
	}
}

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
