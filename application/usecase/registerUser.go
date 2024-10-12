package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/RuanScherer/journey-track-api/adapters/emailtemplateadptr"
	"github.com/RuanScherer/journey-track-api/application/kafka"
	"github.com/RuanScherer/journey-track-api/application/repository"

	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/config"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/matcornic/hermes/v2"
	"gorm.io/gorm"
)

var queueVerificationEmail = doQueueVerificationEmail

type RegisterUserUseCase struct {
	userRepository  repository.UserRepository
	producerFactory kafka.ProducerFactory
}

func NewRegisterUserUseCase(
	userRepository repository.UserRepository,
	producerFactory kafka.ProducerFactory,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepository,
		producerFactory,
	}
}

func (useCase *RegisterUserUseCase) Execute(
	req *appmodel.RegisterUserRequest,
) (*appmodel.RegisterUserResponse, *appmodel.AppError) {
	user, err := model.NewUser(req.Email, req.Name, req.Password)
	if err != nil {
		return nil, appmodel.NewAppError("invalid_data_to_register_user", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.userRepository.Register(user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, appmodel.NewAppError(
				"user_email_already_used",
				"There's already an user using this email",
				appmodel.ErrorTypeValidation,
			)
		}
		return nil, appmodel.NewAppError("unable_to_register_user", "unable to register user", appmodel.ErrorTypeDatabase)
	}

	go queueVerificationEmail(useCase.producerFactory, user)
	return &appmodel.RegisterUserResponse{
		ID:         user.ID,
		Email:      *user.Email,
		Name:       user.Name,
		IsVerified: user.IsVerified,
	}, nil
}

func doQueueVerificationEmail(producerFactory kafka.ProducerFactory, user *model.User) {
	appConfig := config.GetAppConfig()
	verificationLink := fmt.Sprintf(
		"%s/verify-account?userId=%s&token=%s",
		appConfig.FrontendUrl,
		user.ID,
		*user.VerificationToken,
	)
	emailTemplate := hermes.Email{
		Body: hermes.Body{
			Name:  user.Name,
			Title: "Account verification",
			Intros: []string{
				"We are glad you're here!",
				"Please verify your account to start using our services.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to verify your account.",
					Button: hermes.Button{
						Color: "#f25d9c",
						Text:  "Verify account",
						Link:  verificationLink,
					},
				},
			},
			Signature: "Regards",
		},
	}
	content, err := emailtemplateadptr.GenerateEmailHtml(emailTemplate)
	if err != nil {
		slog.Error("Error generating email body as HTML:", "error", err)
		return
	}

	producer, err := producerFactory.NewProducer(map[string]any{
		"bootstrap.servers": appConfig.KafkaBootstrapServers,
		"retries":           3,
		"retry.backoff.ms":  1000,
	})
	if err != nil {
		slog.Error("Error creating kafka producer", "error", err)
		return
	}

	payload, err := json.Marshal(kafka.EmailSendindRequestedPayload{
		To:      *user.Email,
		Subject: "Trackr | Verify your account",
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
