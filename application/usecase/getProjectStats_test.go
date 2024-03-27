package usecase

import (
	"errors"
	"github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestGetProjectStatsUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	projectRepositoryMock := repository.NewMockProjectRepository(ctrl)
	useCase := NewGetProjectStatsUseCase(projectRepositoryMock)

	req := &model.GetProjectStatsRequest{
		ActorID:   "fake-actor-id",
		ProjectID: "fake-project-id",
	}

	projectRepositoryMock.
		EXPECT().
		HasMember(req.ProjectID, req.ActorID).
		Return(false, errors.New("unexpected error"))

	res, err := useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_check_membership] unexpected error")

	projectRepositoryMock.
		EXPECT().
		HasMember(req.ProjectID, req.ActorID).
		Return(false, nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [not_project_member] only project members can see project details")

	projectRepositoryMock.
		EXPECT().
		HasMember(req.ProjectID, req.ActorID).
		AnyTimes().
		Return(true, nil)
	projectRepositoryMock.
		EXPECT().
		FindMembersCountAndEventsCountById(req.ProjectID).
		Return(nil, errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_load_stats] unexpected error")

	projectRepositoryMock.
		EXPECT().
		FindMembersCountAndEventsCountById(req.ProjectID).
		Return(&repository.ProjectInvitesCountAndEventsCount{
			InvitesCount: 10,
			EventsCount:  20,
		}, nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, 10, res.InvitesCount)
	assert.Equal(t, 20, res.EventsCount)
}
