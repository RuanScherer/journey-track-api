package usecase

import (
	"github.com/RuanScherer/journey-track-api/application/jwt"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"golang.org/x/crypto/bcrypt"
)

type SignInUseCase struct {
	userRepository repository.UserRepository
}

func NewSignInUseCase(userRepository repository.UserRepository) *SignInUseCase {
	return &SignInUseCase{userRepository}
}

func (useCase *SignInUseCase) Execute(req *appmodel.SignInRequest) (*appmodel.SignInResponse, *appmodel.AppError) {
	user, err := useCase.userRepository.FindByEmail(req.Email)
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

	jwt, err := jwt.CreateJwtFromUser(user)
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
