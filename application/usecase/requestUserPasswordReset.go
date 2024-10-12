package usecase

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/RuanScherer/journey-track-api/adapters/emailtemplateadptr"
	"github.com/RuanScherer/journey-track-api/application/kafka"
	"github.com/RuanScherer/journey-track-api/application/repository"

	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/config"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/matcornic/hermes/v2"
)

var queuePasswordResetEmail = doQueuePasswordResetEmail

type RequestUserPasswordResetUseCase struct {
	userRepository  repository.UserRepository
	producerFactory kafka.ProducerFactory
}

func NewRequestUserPasswordResetUseCase(
	userRepository repository.UserRepository,
	producerFactory kafka.ProducerFactory,
) *RequestUserPasswordResetUseCase {
	return &RequestUserPasswordResetUseCase{
		userRepository,
		producerFactory,
	}
}

func (useCase *RequestUserPasswordResetUseCase) Execute(req *appmodel.RequestPasswordResetRequest) error {
	user, err := useCase.userRepository.FindByEmail(req.Email)
	if err != nil {
		return appmodel.NewAppError("user_not_found", "user not found", appmodel.ErrorTypeValidation)
	}

	user.RequestPasswordReset()
	err = useCase.userRepository.Save(user)
	if err != nil {
		return appmodel.NewAppError("unable_to_complete", "unable to complete the request", appmodel.ErrorTypeDatabase)
	}

	go queuePasswordResetEmail(useCase.producerFactory, user)
	return nil
}

func doQueuePasswordResetEmail(producerFactory kafka.ProducerFactory, user *model.User) {
	appConfig := config.GetAppConfig()
	passwordResetLink := fmt.Sprintf(
		"%s/reset-password?userId=%s&token=%s",
		appConfig.FrontendUrl,
		user.ID,
		*user.PasswordResetToken,
	)
	emailTemplate := hermes.Email{
		Body: hermes.Body{
			Name:  user.Name,
			Title: "Password reset",
			Intros: []string{
				"It seems you're having trouble with your password.",
				"As requested, we're sending you a link to reset it.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to reset your password.",
					Button: hermes.Button{
						Color: "#f25d9c",
						Text:  "Reset password",
						Link:  passwordResetLink,
					},
				},
			},
			Signature: "Regards",
		},
	}
	content, err := emailtemplateadptr.GenerateEmailHtml(emailTemplate)
	if err != nil {
		slog.Error("Error generating email body as HTML", "error", err)
		return
	}

	producer, err := producerFactory.NewProducer(map[string]any{
		"bootstrap.servers": appConfig.KafkaBootstrapServers,
		"retries":           3,
		"retry.backoff.ms":  1000,
	})
	if err != nil {
		slog.Error("Error setting kafka producer config", "error", err)
		return
	}

	payload, err := json.Marshal(kafka.EmailSendindRequestedPayload{
		To:      *user.Email,
		Subject: "Trackr | Reset your password",
		Content: content,
	})
	if err != nil {
		slog.Error("Error marshalling email sending payload", "error", err)
		return
	}
	message := kafka.Message{Value: payload}
	err = producer.Produce("email-sending-requested", message)
	if err != nil {
		slog.Error("Error producing kafka message to send email", "error", err)
		return
	}
}
