package usecase

import (
	"fmt"
	"log"

	"github.com/RuanScherer/journey-track-api/adapters/email"
	emailutils "github.com/RuanScherer/journey-track-api/adapters/email/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/utils"
	"github.com/RuanScherer/journey-track-api/config"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/RuanScherer/journey-track-api/domain/repository"
	"github.com/matcornic/hermes/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	repository   repository.UserRepository
	emailService email.EmailService
}

func NewUserUseCase(repository repository.UserRepository, emailService email.EmailService) *UserUseCase {
	return &UserUseCase{
		repository,
		emailService,
	}
}

func (useCase *UserUseCase) RegisterUser(req *appmodel.RegisterUserRequest) (*appmodel.RegisterUserResponse, *appmodel.AppError) {
	user, err := model.NewUser(req.Email, req.Name, req.Password)
	if err != nil {
		return nil, appmodel.NewAppError("invalid_data_to_register_user", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.repository.Register(user)
	if err != nil {
		if err == gorm.ErrDuplicatedKey {
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

func (useCase *UserUseCase) sendVerificationEmail(user *model.User) {
	frontendUrl := config.GetAppConfig().FrontendUrl
	verificationLink := fmt.Sprintf(
		"%s/verify-account?userId=%s&token=%s",
		frontendUrl,
		user.ID,
		*user.VerificationToken,
	)
	emailConfig := hermes.Email{
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
	body, err := emailutils.GenerateEmailHtml(emailConfig)
	if err != nil {
		log.Print(err)
		return
	}

	email := email.EmailSendingConfig{
		To:      *user.Email,
		Subject: "Trackr | Verify your account",
		Body:    body,
	}
	err = useCase.emailService.SendEmail(email)
	if err != nil {
		log.Print(err)
		return
	}
}

func (useCase *UserUseCase) VerifyUser(req *appmodel.VerifyUserRequest) *appmodel.AppError {
	user, err := useCase.repository.FindById(req.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appmodel.NewAppError("user_not_found", "user not found", appmodel.ErrorTypeValidation)
		}
		return appmodel.NewAppError("unable_to_find_user", "unable to find user", appmodel.ErrorTypeDatabase)
	}

	err = user.Verify(req.VerificationToken)
	if err != nil {
		if err.Error() == "user already verified" {
			return appmodel.NewAppError("user_already_verified", "user already verified", appmodel.ErrorTypeValidation)
		}
		return appmodel.NewAppError("unable_to_verify_user", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.repository.Save(user)
	if err != nil {
		return appmodel.NewAppError("unable_to_save_user", "unable to save user", appmodel.ErrorTypeDatabase)
	}

	return nil
}

func (useCase *UserUseCase) SignIn(req *appmodel.SignInRequest) (*appmodel.SignInResponse, *appmodel.AppError) {
	user, err := useCase.repository.FindByEmail(req.Email)
	if err != nil {
		return nil, appmodel.NewAppError(
			"invalid_auth_credentials",
			"Invalid authentication credentials",
			appmodel.ErrorTypeValidation,
		)
	}

	if !user.IsVerified {
		return nil, appmodel.NewAppError("user_not_verified", "User is not verified", appmodel.ErrorTypeValidation)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, appmodel.NewAppError(
			"invalid_auth_credentials",
			"Invalid authentication credentials",
			appmodel.ErrorTypeValidation,
		)
	}

	jwt, err := utils.CreateJwtFromUser(user)
	if err != nil {
		return nil, appmodel.NewAppError("unexpected_error", err.Error(), appmodel.ErrorTypeServer)
	}

	return &appmodel.SignInResponse{
		AccessToken: jwt,
		User: appmodel.SignInUser{
			ID:    user.ID,
			Email: *user.Email,
			Name:  user.Name,
		},
	}, nil
}

func (useCase *UserUseCase) EditUser(req *appmodel.EditUserRequest) (*appmodel.EditUserResponse, error) {
	user, err := useCase.repository.FindById(req.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appmodel.NewAppError("user_not_found", "user not found", appmodel.ErrorTypeValidation)
		}
		return nil, appmodel.NewAppError("unable_to_find_user", "unable to find user", appmodel.ErrorTypeDatabase)
	}

	err = user.ChangeName(req.Name)
	if err != nil {
		return nil, appmodel.NewAppError("unable_to_edit_user", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.repository.Save(user)
	if err != nil {
		return nil, appmodel.NewAppError(
			"unable_to_save_user_changes",
			"unable to save user changes",
			appmodel.ErrorTypeDatabase,
		)
	}

	return &appmodel.EditUserResponse{
		ID:         user.ID,
		Email:      *user.Email,
		Name:       user.Name,
		IsVerified: user.IsVerified,
	}, nil
}

func (useCase *UserUseCase) RequestUserPasswordReset(req *appmodel.RequestPasswordResetRequest) error {
	user, err := useCase.repository.FindByEmail(req.Email)
	if err != nil {
		return appmodel.NewAppError("user_not_found", "user not found", appmodel.ErrorTypeValidation)
	}

	user.RequestPasswordReset()
	err = useCase.repository.Save(user)
	if err != nil {
		return appmodel.NewAppError("unable_to_complete", "unable to complete the request", appmodel.ErrorTypeDatabase)
	}

	go useCase.sendPasswordResetEmail(user)
	return nil
}

func (useCase *UserUseCase) sendPasswordResetEmail(user *model.User) {
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

func (useCase *UserUseCase) ResetUserPassword(req *appmodel.PasswordResetRequest) error {
	u, err := useCase.repository.FindById(req.UserID)
	if err != nil {
		return appmodel.NewAppError("user_not_found", "user not found", appmodel.ErrorTypeValidation)
	}

	err = u.ResetPassword(req.Password, req.PasswordResetToken)
	if err != nil {
		if err.Error() == "user has no request for password reset" {
			return appmodel.NewAppError("no_request_for_password_reset", err.Error(), appmodel.ErrorTypeValidation)
		}
		return appmodel.NewAppError("unable_to_reset_password", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.repository.Save(u)
	if err != nil {
		return appmodel.NewAppError(
			"unable_to_save_user_changes",
			"unable to save user changes",
			appmodel.ErrorTypeDatabase,
		)
	}

	return nil
}

func (useCase *UserUseCase) ShowUser(userId string) (*appmodel.ShowUserResponse, error) {
	user, err := useCase.repository.FindById(userId)
	if err != nil {
		return nil, appmodel.NewAppError("user_not_found", "user not found", appmodel.ErrorTypeValidation)
	}

	return &appmodel.ShowUserResponse{
		ID:         user.ID,
		Email:      *user.Email,
		Name:       user.Name,
		IsVerified: user.IsVerified,
	}, nil
}

func (useCase *UserUseCase) SearchUsers(req *appmodel.SearchUsersRequest) (*appmodel.SearchUsersResponse, error) {
	users, err := useCase.repository.Search(repository.UserSearchOptions{
		Email:    req.Email,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	var usersResponse appmodel.SearchUsersResponse
	for _, user := range users {
		usersResponse = append(usersResponse, &appmodel.UserSearchResult{
			ID:    user.ID,
			Email: *user.Email,
			Name:  user.Name,
		})
	}

	return &usersResponse, nil
}
