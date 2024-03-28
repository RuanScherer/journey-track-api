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

func TestDeleteProjectUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	projectRepositoryMock := repository.NewMockProjectRepository(ctrl)
	useCase := NewDeleteProjectUseCase(projectRepositoryMock)

	req := &model.DeleteProjectRequest{
		ActorID:   "fake-actor-id",
		ProjectID: "fake-project-id",
	}

	projectRepositoryMock.
		EXPECT().
		FindById(req.ProjectID).
		Return(nil, gorm.ErrRecordNotFound)

	err := useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [project_not_found] project not found")

	projectRepositoryMock.
		EXPECT().
		FindById(req.ProjectID).
		Return(nil, errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_find_project] unexpected error")

	project, _ := factory.NewProjectWithDefaultOwner("fake project")
	projectRepositoryMock.
		EXPECT().
		FindById(req.ProjectID).
		AnyTimes().
		Return(project, nil)

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [not_project_owner] only project owner can delete project")

	req.ActorID = project.OwnerID
	projectRepositoryMock.
		EXPECT().
		DeleteById(req.ProjectID).
		Return(errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_delete_project] unexpected error")

	projectRepositoryMock.
		EXPECT().
		DeleteById(req.ProjectID).
		Return(nil)

	err = useCase.Execute(req)
	assert.Nil(t, err)
}
