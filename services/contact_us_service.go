package services

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
)

type ContactUsService interface {
	ProcessSubmission(ctx context.Context, request web.ContactUsRequest) error
}
