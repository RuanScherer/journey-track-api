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

func TestCreateProjectUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	projectRepositoryMock := repository.NewMockProjectRepository(ctrl)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	useCase := NewCreateProjectUseCase(projectRepositoryMock, userRepositoryMock)

	req := &model.CreateProjectRequest{
		OwnerID: "fake-owner-id",
		Name:    "",
	}

	userRepositoryMock.
		EXPECT().
		FindById(req.OwnerID).
		Return(nil, gorm.ErrRecordNotFound)

	res, err := useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [project_owner_not_found] project owner not found")

	userRepositoryMock.
		EXPECT().
		FindById(req.OwnerID).
		Return(nil, errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_find_project_owner] unexpected error")

	user, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	userRepositoryMock.
		EXPECT().
		FindById(req.OwnerID).
		AnyTimes().
		Return(user, nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_create_project] project name is required")

	req.Name = "Fake Project"
	projectRepositoryMock.
		EXPECT().
		Register(gomock.Any()).
		Return(errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_save_project] unable to save project")

	projectRepositoryMock.
		EXPECT().
		Register(gomock.Any()).
		Return(nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.ID)
	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, user.ID, res.OwnerID)
}
