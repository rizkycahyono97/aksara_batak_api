package services

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/config"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/utils/helpers"
	"google.golang.org/genai"
	"log/slog"
	"time"
)

type translateServiceImpl struct {
	client *genai.Client
	Log    *slog.Logger
}

func NewTranslateService(log *slog.Logger) TranslateService {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: config.GetEnv("GEMINI_API_KEY", ""),
	})
	if err != nil {
		panic(err)
	}

	return &translateServiceImpl{
		client: client,
		Log:    log,
	}
}

func (c translateServiceImpl) Translate(ctx context.Context, request web.TranslateRequest) (web.TranslateResponse, error) {
	prompt := helpers.BuildTranslationPrompt(request.Text, request.Direction)

	modelConfig := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText("Anda adalah penerjemah ahli Batak-Indonesia. Berikan hanya terjemahan teks tanpa penjelasan tambahan", genai.RoleUser),
	}

	result, err := c.client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		modelConfig,
	)
	if err != nil {
		c.Log.ErrorContext(ctx, "Gemini API error", "error", err)
		return web.TranslateResponse{}, err
	}

	translatedText, err := helpers.ParseResponse(result)
	if err != nil {
		return web.TranslateResponse{}, err
	}

	return web.TranslateResponse{
		OriginalText:   request.Text,
		TranslatedText: translatedText,
		Direction:      request.Direction,
		Timestamp:      time.Now().Format(time.RFC3339),
	}, nil
}
