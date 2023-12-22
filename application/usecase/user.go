package usecase

import (
	"crypto/sha256"
	"encoding/hex"

	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/domain/model"
)

type UserUseCase struct {
	repository model.UserRepository
}

func NewUserUseCase(repository model.UserRepository) *UserUseCase {
	return &UserUseCase{repository: repository}
}

func (useCase *UserUseCase) RegisterUser(req *appmodel.RegisterUserRequest) (*appmodel.RegisterUserResponse, error) {
	passwordHash := sha256.Sum256([]byte(req.Password))
	user, err := model.NewUser(req.Email, req.Name, hex.EncodeToString(passwordHash[:]))
	if err != nil {
		return nil, err
	}

	err = useCase.repository.Register(user)
	if err != nil {
		return nil, err
	}

	return &appmodel.RegisterUserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Name:       user.Name,
		IsVerified: user.IsVerified,
	}, nil
}

func (useCase *UserUseCase) VerifyUser(req *appmodel.VerifyUserRequest) error {
	user, err := useCase.repository.FindById(req.UserID)
	if err != nil {
		return err
	}

	err = user.Verify(req.VerificationToken)
	if err != nil {
		return err
	}

	err = useCase.repository.Save(user)
	if err != nil {
		return err
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
		Email:      user.Email,
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
		Email:      user.Email,
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
			Email:      user.Email,
			Name:       user.Name,
			IsVerified: user.IsVerified,
		})
	}

	return &usersResponse, nil
}
