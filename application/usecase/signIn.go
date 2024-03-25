package usecase

import (
	"github.com/RuanScherer/journey-track-api/application/jwt"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"golang.org/x/crypto/bcrypt"
)

type SignInUseCase struct {
	userRepository repository.UserRepository
	jwtManager     jwt.Manager
}

func NewSignInUseCase(userRepository repository.UserRepository, jwtManager jwt.Manager) *SignInUseCase {
	return &SignInUseCase{userRepository, jwtManager}
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

	token, err := useCase.jwtManager.CreateJwtFromUser(user)
	if err != nil {
		return nil, appmodel.NewAppError("unexpected_error", err.Error(), appmodel.ErrorTypeServer)
	}

	return &appmodel.SignInResponse{
		AccessToken: token,
		User: appmodel.SignInUser{
			ID:    user.ID,
			Email: *user.Email,
			Name:  user.Name,
		},
	}, nil
}
