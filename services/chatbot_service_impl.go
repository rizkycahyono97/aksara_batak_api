package services

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/config"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"log/slog"

	"google.golang.org/genai"
	//"github.com/google/generative-ai-go/genai"
)

type chatbotServiceImpl struct {
	Log *slog.Logger
}

func NewChatbotService(log *slog.Logger) ChatbotService {

	return &chatbotServiceImpl{
		Log: log,
	}
}

func (c chatbotServiceImpl) GeneratePublicResponse(ctx context.Context, request web.ChatbotRequest) (web.ChatbotResponse, error) {
	geminiAPiKey := config.GetEnv("GEMINI_API_KEY", "")

	//client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: geminiAPiKey,
	})
	if err != nil {
		c.Log.InfoContext(ctx, "Gagal membuat client Gemini: %v", err)
	}

	//system
	modelConfig := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(`Kamu adalah seorang ahli dan guru Aksara Batak yang ramah.\n\t\tTugasmu adalah menjawab semua pertanyaan pengguna seputar sejarah, cara penulisan, arti, terjemahan dari bahasa indonesia ke batak maupun sebaliknya, dan contoh penggunaan Aksara Batak.\n\t\t- JAWAB HANYA PERTANYAAN TERKAIT AKSARA BATAK. Jika pengguna bertanya di luar topik, tolak dengan sopan dan arahkan kembali ke topik.\n\t\t- BERIKAN JAWABAN YANG SINGKAT DAN JELAS. Jawabanmu harus ringkas, cukup dalam 2 sampai 4 kalimat saja.`, genai.RoleUser),
	}

	//model gemini
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(request.Message),
		modelConfig,
	)
	if err != nil {
		c.Log.InfoContext(ctx, "Gagal membuat client Gemini: %v", err)
	}

	return web.ChatbotResponse{Reply: result.Text()}, nil
}
