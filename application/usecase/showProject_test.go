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

func TestShowProjectUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	projectRepositoryMock := repository.NewMockProjectRepository(ctrl)
	useCase := NewShowProjectUseCase(projectRepositoryMock, userRepositoryMock)

	req := &model.ShowProjectRequest{
		ActorID:   "fake-actor-id",
		ProjectID: "fake-project-id",
	}

	projectRepositoryMock.
		EXPECT().
		FindById(req.ProjectID).
		Return(nil, gorm.ErrRecordNotFound)

	res, err := useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [project_not_found] project not found")

	projectRepositoryMock.
		EXPECT().
		FindById(req.ProjectID).
		Return(nil, errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_find_project] unexpected error")

	owner, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	project, _ := domainmodel.NewProject("fake-project-id", owner)
	projectRepositoryMock.
		EXPECT().
		FindById(req.ProjectID).
		Return(project, nil)
	userRepositoryMock.
		EXPECT().
		FindById(req.ActorID).
		Return(nil, errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_identify_user] unable to identify the user trying to see project details")

	owner, _ = factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	project, _ = domainmodel.NewProject("fake-project-id", owner)
	actor, _ := domainmodel.NewUser("jane.doe@gmail.com", "Jane Doe", "fake-password")
	projectRepositoryMock.
		EXPECT().
		FindById(req.ProjectID).
		Return(project, nil)
	userRepositoryMock.
		EXPECT().
		FindById(req.ActorID).
		Return(actor, nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [not_project_member] only project members can see project details")

	owner, _ = factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	project, _ = domainmodel.NewProject("fake-project-id", owner)
	req.ActorID = owner.ID
	projectRepositoryMock.
		EXPECT().
		FindById(req.ProjectID).
		Return(project, nil)
	userRepositoryMock.
		EXPECT().
		FindById(req.ActorID).
		Return(owner, nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, project.ID)
	assert.Equal(t, res.Name, project.Name)
	assert.Equal(t, res.OwnerID, project.OwnerID)
	assert.True(t, res.IsOwner)
	assert.Len(t, res.Members, 1)
	assert.Equal(t, res.Members[0].ID, owner.ID)
	assert.Equal(t, res.Members[0].Email, *owner.Email)
	assert.Equal(t, res.Members[0].Name, owner.Name)
}
