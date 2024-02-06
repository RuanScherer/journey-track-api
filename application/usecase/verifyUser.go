package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/domain/repository"
	"gorm.io/gorm"
)

type VerifyUserUseCase struct {
	userRepository repository.UserRepository
}

func NewVerifyUserUseCase(userRepository repository.UserRepository) *VerifyUserUseCase {
	return &VerifyUserUseCase{userRepository}
}

func (useCase *VerifyUserUseCase) Execute(req *appmodel.VerifyUserRequest) *appmodel.AppError {
	user, err := useCase.userRepository.FindById(req.UserID)
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

	err = useCase.userRepository.Save(user)
	if err != nil {
		return appmodel.NewAppError("unable_to_save_user", "unable to save user", appmodel.ErrorTypeDatabase)
	}

	return nil
}
