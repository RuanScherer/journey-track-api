package usecase

import (
	"errors"
	"fmt"
	emailutils "github.com/RuanScherer/journey-track-api/adapters/emailtemplate"
	"github.com/RuanScherer/journey-track-api/application/email"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"log"

	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/config"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/matcornic/hermes/v2"
	"gorm.io/gorm"
)

type RegisterUserUseCase struct {
	userRepository repository.UserRepository
	emailService   email.EmailService
}

func NewRegisterUserUseCase(
	userRepository repository.UserRepository,
	emailService email.EmailService,
) *RegisterUserUseCase {
	return &RegisterUserUseCase{
		userRepository,
		emailService,
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

	go useCase.sendVerificationEmail(user)
	return &appmodel.RegisterUserResponse{
		ID:         user.ID,
		Email:      *user.Email,
		Name:       user.Name,
		IsVerified: user.IsVerified,
	}, nil
}

func (useCase *RegisterUserUseCase) sendVerificationEmail(user *model.User) {
	frontendUrl := config.GetAppConfig().FrontendUrl
	verificationLink := fmt.Sprintf(
		"%s/verify-account?userId=%s&token=%s",
		frontendUrl,
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
	body, err := emailutils.GenerateEmailHtml(emailTemplate)
	if err != nil {
		log.Print(err)
		return
	}

	emailConfig := email.EmailSendingConfig{
		To:      *user.Email,
		Subject: "Trackr | Verify your account",
		Body:    body,
	}
	err = useCase.emailService.SendEmail(emailConfig)
	if err != nil {
		log.Print(err)
		return
	}
}
