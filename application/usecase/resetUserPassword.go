package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/domain/repository"
)

type ResetUserPasswordUseCase struct {
	userRepository repository.UserRepository
}

func NewResetUserPasswordUseCase(userRepository repository.UserRepository) *ResetUserPasswordUseCase {
	return &ResetUserPasswordUseCase{userRepository}
}

func (useCase *ResetUserPasswordUseCase) Execute(req *appmodel.PasswordResetRequest) error {
	u, err := useCase.userRepository.FindById(req.UserID)
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

	err = useCase.userRepository.Save(u)
	if err != nil {
		return appmodel.NewAppError(
			"unable_to_save_user_changes",
			"unable to save user changes",
			appmodel.ErrorTypeDatabase,
		)
	}

	return nil
}
