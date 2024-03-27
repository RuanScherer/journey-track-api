package usecase

import (
	"errors"
	"github.com/RuanScherer/journey-track-api/application/factory"
	"github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
)

func TestEditUserUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	useCase := NewEditUserUseCase(userRepositoryMock)

	req := &model.EditUserRequest{
		UserID: "fake-user-id",
		Name:   "",
	}

	userRepositoryMock.
		EXPECT().
		FindById(req.UserID).
		Return(nil, gorm.ErrRecordNotFound)

	res, err := useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [user_not_found] user not found")

	userRepositoryMock.
		EXPECT().
		FindById(req.UserID).
		Return(nil, errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_find_user] unable to find user")

	user, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	userRepositoryMock.
		EXPECT().
		FindById(req.UserID).
		AnyTimes().
		Return(user, nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_edit_user] name is required")

	req.Name = "Jane Doe"
	userRepositoryMock.
		EXPECT().
		Save(user).
		Return(errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_save_user_changes] unable to save user changes")

	req.Name = "John Doe"
	userRepositoryMock.
		EXPECT().
		Save(user).
		Return(nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, user.ID, res.ID)
	assert.Equal(t, *user.Email, res.Email)
	assert.Equal(t, user.Name, res.Name)
	assert.Equal(t, user.IsVerified, res.IsVerified)
}
