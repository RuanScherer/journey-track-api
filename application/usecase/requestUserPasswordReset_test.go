package usecase

import (
	"errors"
	"github.com/RuanScherer/journey-track-api/application/email"
	"github.com/RuanScherer/journey-track-api/application/factory"
	"github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRequestUserPasswordResetUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	emailServiceMock := email.NewMockEmailService(ctrl)
	useCase := NewRequestUserPasswordResetUseCase(userRepositoryMock, emailServiceMock)

	req := &model.RequestPasswordResetRequest{
		Email: "john.doe@gmail.com",
	}

	userRepositoryMock.
		EXPECT().
		FindByEmail(req.Email).
		Return(nil, errors.New("user not found"))

	err := useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [user_not_found] user not found")

	user, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	userRepositoryMock.
		EXPECT().
		FindByEmail(req.Email).
		AnyTimes().
		Return(user, nil)
	userRepositoryMock.
		EXPECT().
		Save(user).
		Return(errors.New("unable to save"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_complete] unable to complete the request")

	userRepositoryMock.
		EXPECT().
		Save(user).
		AnyTimes().
		Return(nil)
	emailServiceMock.
		EXPECT().
		SendEmail(gomock.Any()).
		AnyTimes()

	err = useCase.Execute(req)
	assert.Nil(t, err)
}

func TestRequestUserPasswordResetUseCase_sendPasswordResetEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	emailServiceMock := email.NewMockEmailService(ctrl)
	useCase := NewRequestUserPasswordResetUseCase(userRepositoryMock, emailServiceMock)

	user, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	user.RequestPasswordReset()

	emailServiceMock.
		EXPECT().
		SendEmail(gomock.Any()).
		Return(errors.New("unable to send email"))
	useCase.sendPasswordResetEmail(user)
}
