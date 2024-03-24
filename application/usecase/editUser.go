package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"gorm.io/gorm"
)

type EditUserUseCase struct {
	userRepository repository.UserRepository
}

func NewEditUserUseCase(userRepository repository.UserRepository) *EditUserUseCase {
	return &EditUserUseCase{userRepository}
}

func (useCase *EditUserUseCase) Execute(req *appmodel.EditUserRequest) (*appmodel.EditUserResponse, error) {
	user, err := useCase.userRepository.FindById(req.UserID)
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

	err = useCase.userRepository.Save(user)
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
