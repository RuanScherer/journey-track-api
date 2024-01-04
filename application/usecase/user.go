package usecase

import (
	"crypto/sha256"
	"encoding/hex"

	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"gorm.io/gorm"
)

type UserUseCase struct {
	repository model.UserRepository
}

func NewUserUseCase(repository model.UserRepository) *UserUseCase {
	return &UserUseCase{repository: repository}
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
				appmodel.ErrorTypeDatabase,
			)
		}
		return nil, appmodel.NewAppError("unable_to_register_user", "unable to register user", appmodel.ErrorTypeDatabase)
	}

	return &appmodel.RegisterUserResponse{
		ID:         user.ID,
		Email:      *user.Email,
		Name:       user.Name,
		IsVerified: user.IsVerified,
	}, nil
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
		return appmodel.NewAppError("unable_to_verify_user", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.repository.Save(user)
	if err != nil {
		return appmodel.NewAppError("unable_to_save_user", "unable to save user", appmodel.ErrorTypeDatabase)
	}

	return nil
}

func (useCase *UserUseCase) EditUser(req *appmodel.EditUserRequest) (*appmodel.EditUserResponse, error) {
	user, err := useCase.repository.FindById(req.UserID)
	if err != nil {
		return nil, err
	}

	err = user.ChangeName(req.Name)
	if err != nil {
		return nil, err
	}

	err = useCase.repository.Save(user)
	if err != nil {
		return nil, err
	}

	return &appmodel.EditUserResponse{
		ID:         user.ID,
		Email:      *user.Email,
		Name:       user.Name,
		IsVerified: user.IsVerified,
	}, nil
}

func (useCase *UserUseCase) RequestUserPasswordReset(userId string) error {
	user, err := useCase.repository.FindById(userId)
	if err != nil {
		return err
	}

	user.RequestPasswordReset()
	err = useCase.repository.Save(user)
	if err != nil {
		return err
	}

	return nil
}

func (useCase *UserUseCase) ResetUserPassword(req *appmodel.PasswordResetRequest) error {
	u, err := useCase.repository.FindById(req.UserID)
	if err != nil {
		return err
	}

	passwordHash := sha256.Sum256([]byte(req.Password))
	err = u.ResetPassword(hex.EncodeToString(passwordHash[:]), req.PasswordResetToken)
	if err != nil {
		return err
	}

	err = useCase.repository.Save(u)
	if err != nil {
		return err
	}

	return nil
}

func (useCase *UserUseCase) ShowUser(userId string) (*appmodel.ShowUserResponse, error) {
	user, err := useCase.repository.FindById(userId)
	if err != nil {
		return nil, err
	}

	return &appmodel.ShowUserResponse{
		ID:         user.ID,
		Email:      *user.Email,
		Name:       user.Name,
		IsVerified: user.IsVerified,
	}, nil
}

func (useCase *UserUseCase) SearchUsersByEmail(email string) (*appmodel.SearchUserResponse, error) {
	users, err := useCase.repository.SearchByEmail(email)
	if err != nil {
		return nil, err
	}

	var usersResponse appmodel.SearchUserResponse
	for _, user := range users {
		usersResponse = append(usersResponse, &appmodel.UserSearchResult{
			ID:         user.ID,
			Email:      *user.Email,
			Name:       user.Name,
			IsVerified: user.IsVerified,
		})
	}

	return &usersResponse, nil
}
