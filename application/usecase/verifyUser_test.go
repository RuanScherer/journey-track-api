package usecase

import (
	"errors"
	"github.com/RuanScherer/journey-track-api/application/factory"
	"github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	domainmodel "github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
)

func TestVerifyUserUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	useCase := NewVerifyUserUseCase(userRepositoryMock)

	req := &model.VerifyUserRequest{
		UserID:            "fake-user-id",
		VerificationToken: "fake-verification-token",
	}

	userRepositoryMock.
		EXPECT().
		FindById(req.UserID).
		Return(nil, gorm.ErrRecordNotFound)

	err := useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [user_not_found] user not found")

	userRepositoryMock.
		EXPECT().
		FindById(req.UserID).
		Return(nil, errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_find_user] unable to find user")

	user, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	userRepositoryMock.
		EXPECT().
		FindById(req.UserID).
		Return(user, nil)

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [user_already_verified] user already verified")

	user, _ = domainmodel.NewUser("john.doe@gmail.com", "John Doe", "fake-password")
	userRepositoryMock.
		EXPECT().
		FindById(req.UserID).
		Return(user, nil)

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_verify_user] invalid verification token")

	user, _ = domainmodel.NewUser("john.doe@gmail.com", "John Doe", "fake-password")
	req.VerificationToken = *user.VerificationToken
	userRepositoryMock.
		EXPECT().
		FindById(req.UserID).
		Return(user, nil)
	userRepositoryMock.
		EXPECT().
		Save(gomock.Any()).
		Return(errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_save_user] unable to save user")

	user, _ = domainmodel.NewUser("john.doe@gmail.com", "John Doe", "fake-password")
	req.VerificationToken = *user.VerificationToken
	userRepositoryMock.
		EXPECT().
		FindById(req.UserID).
		Return(user, nil)
	userRepositoryMock.
		EXPECT().
		Save(user).
		Return(nil)

	err = useCase.Execute(req)
	assert.Nil(t, err)
}
