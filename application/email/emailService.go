package email

type EmailService interface {
	SendEmail(email EmailSendingConfig) error
}

type EmailSendingConfig struct {
	To      string `validate:"required~Receiver email is required,email~Receiver email is invalid"`
	Subject string `validate:"required~Subject is required,notblank~Subject must not be blank"`
	Body    string `validate:"required~Body is required,notblank~Body must not be blank"`
}
