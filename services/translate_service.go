package services

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
)

type TranslateService interface {
	Translate(ctx context.Context, request web.TranslateRequest) (web.TranslateResponse, error)
}
