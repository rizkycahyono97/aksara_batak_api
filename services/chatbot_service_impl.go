package services

import (
	"context"
	"fmt"
	"github.com/rizkycahyono97/aksara_batak_api/config"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"google.golang.org/api/option"
	"log/slog"
	"strings"

	//"google.golang.org/genai"
	"github.com/google/generative-ai-go/genai"
)

type chatbotServiceImpl struct {
	GeminiClient *genai.GenerativeModel
	Log          *slog.Logger
}

func NewChatbotService(log *slog.Logger) ChatbotService {
	geminiAPiKey := config.GetEnv("GEMINI_API_KEY", "")
	ctx := context.Background()

	//client
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiAPiKey))
	if err != nil {
		log.InfoContext(ctx, "Gagal membuat client Gemini: %v", err)
	}
	defer client.Close()

	//model gemini
	model := client.GenerativeModel("gemini-2.5-flash")

	return &chatbotServiceImpl{
		GeminiClient: model,
		Log:          log,
	}
}

func (c chatbotServiceImpl) GeneratePublicResponse(ctx context.Context, request web.ChatbotRequest) (web.ChatbotResponse, error) {
	//system promt
	systemPromt := `
		Kamu adalah seorang ahli dan guru Aksara Batak yang ramah.
		Tugasmu adalah menjawab semua pertanyaan pengguna seputar sejarah, cara penulisan, arti, dan contoh penggunaan Aksara Batak.
		- JAWAB HANYA PERTANYAAN TERKAIT AKSARA BATAK. Jika pengguna bertanya di luar topik, tolak dengan sopan dan arahkan kembali ke topik.
		- BERIKAN JAWABAN YANG SINGKAT DAN JELAS. Jawabanmu harus ringkas, cukup dalam 4-6 kalimat saja.
	`

	//menggabungkan systempromt dengna pertanyaan user
	fullPromt := fmt.Sprintf("%s\n\nPertanyaan Pengguna: %s", systemPromt, request.Message)

	//penggil geminiapi
	resp, err := c.GeminiClient.GenerateContent(ctx, genai.Text(fullPromt))
	if err != nil {
		c.Log.InfoContext(ctx, "Error memanggil Gemini API: %v", err)
		return web.ChatbotResponse{}, err
	}

	//ekstrak text dari response
	var botReply strings.Builder
	if resp != nil && len(resp.Candidates) > 0 {
		for _, part := range resp.Candidates[0].Content.Parts {
			if txt, ok := part.(genai.Text); ok {
				botReply.WriteString(string(txt))
			}
		}
	}
	if botReply.Len() == 0 {
		c.Log.InfoContext(ctx, "Nothing to send", "response", resp)
		return web.ChatbotResponse{}, err
	}

	return web.ChatbotResponse{Reply: botReply.String()}, nil
}
