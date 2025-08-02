package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
)

type ContactRepository interface {
	Create(ctx context.Context, submission domain.ContactSubmissions) (domain.ContactSubmissions, error)
}
