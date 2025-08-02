package services

import "context"

type EmailService interface {
	SendContactNotification(ctx context.Context, subject string, body string) error
}
