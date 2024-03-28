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

func TestEditProjectUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	projectRepositoryMock := repository.NewMockProjectRepository(ctrl)
	useCase := NewEditProjectUseCase(projectRepositoryMock)

	req := &model.EditProjectRequest{
		ActorID:   "fake-actor-id",
		ProjectID: "fake-project-id",
		Name:      "",
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

	project, _ := factory.NewProjectWithDefaultOwner("fake proejct")
	projectRepositoryMock.
		EXPECT().
		FindById(req.ProjectID).
		AnyTimes().
		Return(project, nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [not_project_owner] only project owner can edit project")

	req.ActorID = project.OwnerID
	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_edit_project] project name cannot be empty")

	req.Name = "edited fake project"
	projectRepositoryMock.
		EXPECT().
		Save(project).
		Return(errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_save_project_changes] unable to save project changes")

	projectRepositoryMock.
		EXPECT().
		Save(project).
		Return(nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, project.ID, res.ID)
	assert.Equal(t, req.Name, res.Name)
	assert.Equal(t, project.OwnerID, res.OwnerID)
}
