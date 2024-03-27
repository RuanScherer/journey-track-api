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

func TestListProjectInvitesUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	projectInviteMockRepository := repository.NewMockProjectInviteRepository(ctrl)
	projectMockRepository := repository.NewMockProjectRepository(ctrl)
	useCase := NewListProjectInvitesUseCase(projectInviteMockRepository, projectMockRepository)

	req := &model.ListProjectInvitesRequest{
		ActorID:   "fake-actor-id",
		ProjectID: "fake-project-id",
		Status:    "",
	}

	projectMockRepository.
		EXPECT().
		HasMember(req.ProjectID, req.ActorID).
		Return(false, errors.New("unexpected error"))

	res, err := useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_check_membership] unexpected errpr")

	projectMockRepository.
		EXPECT().
		HasMember(req.ProjectID, req.ActorID).
		Return(false, nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [not_project_member] only project members can see the project invites")

	projectMockRepository.
		EXPECT().
		HasMember(req.ProjectID, req.ActorID).
		AnyTimes().
		Return(true, nil)
	projectInviteMockRepository.
		EXPECT().
		ListByProjectAndStatus(req.ProjectID, domainmodel.ProjectInviteStatusPending).
		Return(nil, errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_find_project_invites] unable to find the project invites")

	projectInviteMockRepository.
		EXPECT().
		ListByProjectAndStatus(req.ProjectID, domainmodel.ProjectInviteStatusPending).
		Return(nil, gorm.ErrRecordNotFound)

	res, err = useCase.Execute(req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Empty(t, *res)

	req.Status = domainmodel.ProjectInviteStatusAccepted
	user1, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	user2, _ := factory.NewVerifiedUser("jane.doe@gmail.com", "Jane Doe", "fake-password")
	project, _ := factory.NewProjectWithDefaultOwner("fake project")
	invitation1, _ := domainmodel.NewProjectInvite(project, user1)
	invitation2, _ := domainmodel.NewProjectInvite(project, user2)
	projectInviteMockRepository.
		EXPECT().
		ListByProjectAndStatus(req.ProjectID, domainmodel.ProjectInviteStatusAccepted).
		Return([]*domainmodel.ProjectInvite{invitation1, invitation2}, nil)

	res, err = useCase.Execute(req)
	invitations := *res
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Len(t, invitations, 2)
	assert.Equal(t, invitations[0].ID, invitation1.ID)
	assert.Equal(t, invitations[0].Project.ID, invitation1.Project.ID)
	assert.Equal(t, invitations[0].User.ID, invitation1.User.ID)
	assert.Equal(t, invitations[0].Status, invitation1.Status)
	assert.Equal(t, invitations[1].ID, invitation2.ID)
	assert.Equal(t, invitations[1].Project.ID, invitation2.Project.ID)
	assert.Equal(t, invitations[1].User.ID, invitation2.User.ID)
	assert.Equal(t, invitations[1].Status, invitation2.Status)
}
