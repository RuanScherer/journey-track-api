package usecase

import (
	"errors"
	"github.com/RuanScherer/journey-track-api/application/repository"
	domainmodel "github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestShowUserUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	useCase := NewShowUserUseCase(userRepositoryMock)
	userId := "fake-user-id"

	userRepositoryMock.
		EXPECT().
		FindById(userId).
		Return(nil, errors.New("user not found"))

	res, err := useCase.Execute(userId)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [user_not_found] user not found")

	user, _ := domainmodel.NewUser("john.doe@gmail.com", "John Doe", "fake-password")
	userRepositoryMock.
		EXPECT().
		FindById(userId).
		Return(user, nil)

	res, err = useCase.Execute(userId)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, user.ID)
	assert.Equal(t, res.Email, *user.Email)
	assert.Equal(t, res.Name, user.Name)
	assert.Equal(t, res.IsVerified, user.IsVerified)
}
