package usecase

import (
	"errors"
	"github.com/RuanScherer/journey-track-api/application/factory"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestResetUserPasswordUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	useCase := NewResetUserPasswordUseCase(userRepositoryMock)

	req := &appmodel.PasswordResetRequest{
		UserID:             "fake-user-id",
		Password:           "fake-password",
		PasswordResetToken: "fake-password-reset-token",
	}

	userRepositoryMock.
		EXPECT().
		FindById(req.UserID).
		Return(nil, errors.New("user not found"))

	err := useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [user_not_found] user not found")

	user, _ := factory.NewVerifiedUser("john.doe@gmil.com", "John Doe", "fake-password")
	userRepositoryMock.
		EXPECT().
		FindById(req.UserID).
		Return(user, nil)

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [no_request_for_password_reset] user has no request for password reset")

	user.RequestPasswordReset()
	userRepositoryMock.
		EXPECT().
		FindById(req.UserID).
		AnyTimes().
		Return(user, nil)

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_reset_password] invalid password reset token")

	req.PasswordResetToken = *user.PasswordResetToken
	userRepositoryMock.
		EXPECT().
		Save(user).
		Return(errors.New("unable to save user changes"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_save_user_changes] unable to save user changes")

	user.RequestPasswordReset()
	req.PasswordResetToken = *user.PasswordResetToken
	userRepositoryMock.
		EXPECT().
		Save(user).
		Return(nil)

	err = useCase.Execute(req)
	assert.Nil(t, err)
}
