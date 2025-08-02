package repositories

import (
	"context"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
)

type ContactSubmissionsRepository interface {
	Create(ctx context.Context, submission domain.ContactSubmissions) (domain.ContactSubmissions, error)
}
