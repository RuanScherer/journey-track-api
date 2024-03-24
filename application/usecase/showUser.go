package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
)

type ShowUserUseCase struct {
	userRepository repository.UserRepository
}

func NewShowUserUseCase(userRepository repository.UserRepository) *ShowUserUseCase {
	return &ShowUserUseCase{userRepository}
}

func (useCase *ShowUserUseCase) Execute(userId string) (*appmodel.ShowUserResponse, error) {
	user, err := useCase.userRepository.FindById(userId)
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
