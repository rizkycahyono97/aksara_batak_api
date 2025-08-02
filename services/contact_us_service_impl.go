package services

import (
	"context"
	"fmt"
	"github.com/rizkycahyono97/aksara_batak_api/model/domain"
	"github.com/rizkycahyono97/aksara_batak_api/model/web"
	"github.com/rizkycahyono97/aksara_batak_api/repositories"
	"log/slog"
)

type contactUsServiceImpl struct {
	ContactRepo  repositories.ContactRepository
	EmailService EmailService
	Log          *slog.Logger
}

func NewContactUsService(
	contactRepo repositories.ContactRepository,
	emailService EmailService,
	log *slog.Logger,
) ContactUsService {
	return &contactUsServiceImpl{
		ContactRepo:  contactRepo,
		EmailService: emailService,
		Log:          log,
	}
}

func (c contactUsServiceImpl) ProcessSubmission(ctx context.Context, request web.ContactUsRequest) error {
	c.Log.Info("Memulai pemrosesan pesan 'Contact Us'", "from_email", request.Email)

	submission := domain.ContactSubmissions{
		Name:    request.Name,
		Email:   request.Email,
		Message: request.Message,
		Status:  "baru",
	}

	_, err := c.ContactRepo.Create(ctx, submission)
	if err != nil {
		c.Log.Error("Gagal menyimpan pesan 'Contact Us' ke database", "error", err)
		return fmt.Errorf("gagal menyimpan pesan Anda ke database")
	}
	c.Log.Info("Pesan berhasil disimpan ke database", "from_email", request.Email)

	subject := fmt.Sprintf("Pesan Baru 'Poda-Horas' dari: %s", request.Name)
	body := fmt.Sprintf(
		"Anda telah menerima pesan baru melalui formulir 'Contact-Us Poda Horas'.\n\n"+
			"Dari: %s\n"+
			"Email: %s\n"+
			"Message: %s\n",
		request.Name, request.Email, request.Message,
	)

	err = c.EmailService.SendContactNotification(ctx, subject, body)
	if err != nil {
		c.Log.Error(
			"PROSES SUKSES NAMUN GAGAL MENGIRIM NOTIFIKASI EMAIL",
			"from_email", request.Email,
			"error", err,
		)
	}

	c.Log.Info("Pemrosesan pesan 'Contact Us' selesai", "from_email", request.Email)
	return nil
}
