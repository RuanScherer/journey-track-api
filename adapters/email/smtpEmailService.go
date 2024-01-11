package email

import (
	"github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/config"
	"github.com/asaskevich/govalidator"
	"gopkg.in/gomail.v2"
)

type SmtpEmailService struct {
	dialer *gomail.Dialer
}

func NewSmtpEmailService() *SmtpEmailService {
	appConfig := config.GetAppConfig()

	return &SmtpEmailService{
		dialer: gomail.NewDialer(
			appConfig.EmailSmtpHost,
			int(appConfig.EmailSmtpPort),
			appConfig.EmailSmtpUsername,
			appConfig.EmailSmtpPassword,
		),
	}
}

func (service *SmtpEmailService) SendEmail(email EmailSendingConfig) error {
	_, err := govalidator.ValidateStruct(&email)
	if err != nil {
		return model.NewAppError("invalid_email_sending_config", err.Error(), model.ErrorTypeValidation)
	}

	message := gomail.NewMessage()
	message.SetHeader("From", service.dialer.Username)
	message.SetHeader("To", email.To)
	message.SetHeader("Subject", email.Subject)
	message.SetBody("text/html", email.Body)

	err = service.dialer.DialAndSend(message)
	if err != nil {
		return model.NewAppError("unable_to_send_email", err.Error(), model.ErrorTypeServer)
	}
	return nil
}
