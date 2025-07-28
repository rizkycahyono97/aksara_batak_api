package services

import "C"
import (
	"context"
	"fmt"
	"github.com/rizkycahyono97/aksara_batak_api/config"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/repositories"
	"github.com/rizkycahyono97/aksara_batak_api/utils/helpers"
	"log/slog"
	"time"

	"google.golang.org/genai"
	//"github.com/google/generative-ai-go/genai"
)

type chatbotServiceImpl struct {
	GeminiClient   *genai.Client
	ChatRepository repositories.ChatHistoryRepository
	Log            *slog.Logger
	ModelName      string
}

func NewChatbotService(
	chatRepository repositories.ChatHistoryRepository,
	log *slog.Logger,
) ChatbotService {

	return &chatbotServiceImpl{
		ChatRepository: chatRepository,
		Log:            log,
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

func (c chatbotServiceImpl) GeneratePrivateResponse(ctx context.Context, request web.ChatbotRequest) (web.ChatbotResponse, error) {
	c.Log.InfoContext(ctx, "Processing chat message", "user_id", request.Userid, "message_length", len(request.Message))

	geminiAPiKey := config.GetEnv("GEMINI_API_KEY", "")

	//client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: geminiAPiKey,
	})
	if err != nil {
		c.Log.InfoContext(ctx, "Gagal membuat client Gemini: %v", err)
	}

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

	//ambil history
	start := time.Now()
	histories, err := c.ChatRepository.GetLastFiveByUserID(ctx, request.Userid)
	if err != nil {
		c.Log.ErrorContext(ctx, "Failed to get chat history",
			"error", err,
			"duration_ms", time.Since(start).Milliseconds(),
		)
		return web.ChatbotResponse{}, fmt.Errorf("Failed to get chat history: %w", err)
	}
	c.Log.DebugContext(ctx, "Retrieved chat history", "count", len(histories), "duration_ms", time.Since(start).Milliseconds())

	var historyResponse []web.ChatHistories
	for _, h := range histories {
		historyResponse = append(historyResponse, web.ChatHistories{
			ID:        int64(h.ID),
			UserID:    h.UserID,
			Message:   h.Message,
			Reply:     h.Reply,
			CreatedAt: h.CreatedAt,
			UpdatedAt: h.UpdatedAt,
		})
	}

	prompt := helpers.BuildPrompt(request.Message, histories)
	c.Log.DebugContext(ctx, "Generated prompt", "prompt", len(prompt))

	//gemini
	startGen := time.Now()
	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		modelConfig,
	)
	if err != nil {
		c.Log.ErrorContext(ctx, "Gemini API failed",
			"error", err,
			"duration_ms", time.Since(startGen).Milliseconds(),
		)
		return web.ChatbotResponse{}, fmt.Errorf("gemini API error: %w", err)
	}

	reply, err := helpers.ParseResponse(resp)
	if err != nil {
		c.Log.ErrorContext(ctx, "Failed to parse Gemini response",
			"error", err,
			"raw_response", fmt.Sprintf("%+v", resp),
		)
		return web.ChatbotResponse{}, fmt.Errorf("parse error: %w", err)
	}

	//simpan history
	newChat := &domain.ChatHistories{
		UserID:    request.Userid,
		Message:   request.Message,
		Reply:     reply,
		CreatedAt: time.Now(),
	}
	if err := c.ChatRepository.Create(ctx, newChat); err != nil {
		c.Log.ErrorContext(ctx, "Failed to save chat history", "error", err)
	}

	c.Log.InfoContext(ctx, "Successfully processed message",
		"user_id", request.Userid,
		"reply_length", len(reply),
		"total_duration_ms", time.Since(start).Milliseconds(),
	)

	return web.ChatbotResponse{
		Reply:   reply,
		History: historyResponse,
	}, nil
}
