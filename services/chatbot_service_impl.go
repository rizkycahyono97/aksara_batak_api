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
		SystemInstruction: genai.NewContentFromText(`
			Kamu adalah 'Dongan Digital', seorang ahli budaya Batak yang ramah dan informatif.
			
			Tugas utamamu adalah menjawab segala pertanyaan tentang budaya Batak secara akurat dan mudah dimengerti. Cakupan pengetahuanmu meliputi:
			- AKSARA BATAK: Sejarah, cara penulisan, arti filosofis, dan penggunaannya.
			- BAHASA BATAK: Terjemahan dari/ke Bahasa Indonesia, kosakata, dan ungkapan umum.
			- SASTRA & PUISI: Arti dan contoh dari umpasa, pantun, dan peribahasa Batak.
			- ADAT & TRADISI: Penjelasan tentang sistem marga, pernikahan, upacara adat, dan silsilah (tarombo).
			- SEJARAH & KULINER: Sejarah suku Batak dan penjelasan tentang makanan khas Batak.
			
			Aturan Wajib:
			1. JAWAB SINGKAT: Berikan jawaban yang singkat, jelas, dan padat dalam 2-4 kalimat.
			2. FOKUS PADA TOPIK: Jika pertanyaan di luar topik budaya Batak, tolak dengan sopan dan kembalikan percakapan ke topik Batak.
			3. BAHASA: Gunakan Bahasa Indonesia yang baik kecuali pengguna meminta dalam bahasa lain.`,
			genai.RoleUser),
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
