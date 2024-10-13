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

func TestTrackEventUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockProjectRepository := repository.NewMockProjectRepository(ctrl)
	mockEventRepository := repository.NewMockEventRepository(ctrl)
	useCase := NewTrackEventUseCase(mockProjectRepository, mockEventRepository)

	req := &appmodel.TrackEventRequest{Name: "fake event", ProjectToken: "fake-project-token"}

	mockProjectRepository.
		EXPECT().
		FindByToken(req.ProjectToken).
		Return(nil, errors.New("unexpected error"))

	err := useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "project not found")

	req.Name = ""
	project, _ := factory.NewProjectWithDefaultOwner("fake project")
	mockProjectRepository.
		EXPECT().
		FindByToken(req.ProjectToken).
		Return(project, nil)

	err = useCase.Execute(req)
	assert.NotNil(t, err)

	req.Name = "fake event"
	mockProjectRepository.
		EXPECT().
		FindByToken(req.ProjectToken).
		Return(project, nil)
	mockEventRepository.
		EXPECT().
		Register(gomock.Any()).
		Return(errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "unexpected error")

	mockProjectRepository.
		EXPECT().
		FindByToken(req.ProjectToken).
		Return(project, nil)
	mockEventRepository.
		EXPECT().
		Register(gomock.Any()).
		Return(nil)

	err = useCase.Execute(req)
	assert.Nil(t, err)
}
