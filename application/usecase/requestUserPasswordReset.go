package usecase

import (
	"fmt"
	"log"

	"github.com/RuanScherer/journey-track-api/adapters/email"
	emailutils "github.com/RuanScherer/journey-track-api/adapters/email/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/config"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/RuanScherer/journey-track-api/domain/repository"
	"github.com/matcornic/hermes/v2"
)

type RequestUserPasswordResetUseCase struct {
	userRepository repository.UserRepository
	emailService   email.EmailService
}

func NewRequestUserPasswordResetUseCase(
	userRepository repository.UserRepository,
	emailService email.EmailService,
) *RequestUserPasswordResetUseCase {
	return &RequestUserPasswordResetUseCase{
		userRepository,
		emailService,
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

	go useCase.sendPasswordResetEmail(user)
	return nil
}

func (useCase *RequestUserPasswordResetUseCase) sendPasswordResetEmail(user *model.User) {
	frontendUrl := config.GetAppConfig().FrontendUrl
	passwordResetLink := fmt.Sprintf(
		"%s/reset-password?userId=%s&token=%s",
		frontendUrl,
		user.ID,
		*user.PasswordResetToken,
	)
	emailConfig := hermes.Email{
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
	body, err := emailutils.GenerateEmailHtml(emailConfig)
	if err != nil {
		log.Print(err)
		return
	}

	email := email.EmailSendingConfig{
		To:      *user.Email,
		Subject: "Trackr | Reset your password",
		Body:    body,
	}
	err = useCase.emailService.SendEmail(email)
	if err != nil {
		log.Print(err)
		return
	}
}
