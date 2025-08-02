package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/smtp"
	"os"
)

type emailServiceImpl struct {
	Log *slog.Logger
}

func NewEmailService(log *slog.Logger) EmailService {
	return &emailServiceImpl{
		Log: log,
	}
}

func (e emailServiceImpl) SendContactNotification(ctx context.Context, subject string, body string) error {
	if subject == "" || body == "" {
		err := errors.New("subject or body is empty")
		e.Log.Warn("Validasi email gagal", "error", err)
		return err
	}

	//env
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_SENDER_EMAIL")
	senderPass := os.Getenv("SMTP_SENDER_PASSWORD")
	receiverEmail := os.Getenv("SMTP_RECEIVER_EMAIL")

	if host == "" || port == "" || senderEmail == "" || senderPass == "" || receiverEmail == "" {
		err := errors.New("konfigurasi SMPT .env tidak lengkap")
		e.Log.Error("Gagal mengirim email karena konfigurasi tidak lengkap", "error", err)
		return err
	}
	e.Log.Info("Mencoba mengirim notifikasi email...", "to", receiverEmail, "subject", subject)

	auth := smtp.PlainAuth("", senderEmail, senderPass, host)

	//format email
	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n", receiverEmail, subject, body))

	addr := fmt.Sprintf("%s:%s", port, host)

	err := smtp.SendMail(addr, auth, senderEmail, []string{receiverEmail}, message)
	if err != nil {
		e.Log.Error("Gagal mengirim email via SMTP", "error", err)
		return fmt.Errorf("gagal mengirim email: %w", err)
	}

	e.Log.Info("Email notifikasi berhasil terkirim", "to", receiverEmail)
	return nil
}
