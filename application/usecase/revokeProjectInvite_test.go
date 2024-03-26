package usecase

import (
	"errors"
	"github.com/RuanScherer/journey-track-api/application/factory"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	domainmodel "github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
)

func TestRevokeProjectInviteUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	projectInviteMockRepository := repository.NewMockProjectInviteRepository(ctrl)
	userMockRepository := repository.NewMockUserRepository(ctrl)
	useCase := NewRevokeProjectInviteUseCase(projectInviteMockRepository, userMockRepository)

	req := &appmodel.RevokeProjectInviteRequest{
		ActorID:         "fake-actor-id",
		ProjectInviteID: "fake-project-invite-id",
	}

	projectInviteMockRepository.
		EXPECT().
		FindById(req.ProjectInviteID).
		Return(nil, gorm.ErrRecordNotFound)

	err := useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [project_invite_not_found] project invite not found")

	projectInviteMockRepository.
		EXPECT().
		FindById(req.ProjectInviteID).
		Return(nil, errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_find_project_invite] unexpected error")

	project, _ := factory.NewProjectWithDefaultOwner("fake project")
	user, _ := factory.NewVerifiedUser("jane.doe@gmail.com", "Jane Doe", "fake-password")
	invitation, _ := domainmodel.NewProjectInvite(project, user)
	_ = invitation.Accept(*invitation.Token)
	projectInviteMockRepository.
		EXPECT().
		FindById(req.ProjectInviteID).
		Return(invitation, nil)

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_revoke_project_invite] invite already answered or revoked")

	project, _ = factory.NewProjectWithDefaultOwner("fake project")
	user, _ = factory.NewVerifiedUser("jane.doe@gmail.com", "Jane Doe", "fake-password")
	invitation, _ = domainmodel.NewProjectInvite(project, user)
	projectInviteMockRepository.
		EXPECT().
		FindById(req.ProjectInviteID).
		Return(invitation, nil)
	userMockRepository.
		EXPECT().
		FindById(req.ActorID).
		Return(nil, errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_identify_user] unable to identify the user trying to revoke the invite")

	project, _ = factory.NewProjectWithDefaultOwner("fake project")
	user, _ = factory.NewVerifiedUser("jane.doe@gmail.com", "Jane Doe", "fake-password")
	invitation, _ = domainmodel.NewProjectInvite(project, user)
	randomUser, _ := factory.NewVerifiedUser("random@gmail.com", "random user", "fake-password")
	projectInviteMockRepository.
		EXPECT().
		FindById(req.ProjectInviteID).
		Return(invitation, nil)
	userMockRepository.
		EXPECT().
		FindById(req.ActorID).
		Return(randomUser, nil)

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [not_project_member] only project members can revoke invites")

	project, _ = factory.NewProjectWithDefaultOwner("fake project")
	user, _ = factory.NewVerifiedUser("jane.doe@gmail.com", "Jane Doe", "fake-password")
	invitation, _ = domainmodel.NewProjectInvite(project, user)
	req.ActorID = project.Members[0].ID
	projectInviteMockRepository.
		EXPECT().
		FindById(req.ProjectInviteID).
		Return(invitation, nil)
	userMockRepository.
		EXPECT().
		FindById(project.Members[0].ID).
		Return(project.Members[0], nil)
	projectInviteMockRepository.
		EXPECT().
		DeleteById(req.ProjectInviteID).
		Return(errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_delete_project_invite] unexpected error")

	project, _ = factory.NewProjectWithDefaultOwner("fake project")
	user, _ = factory.NewVerifiedUser("jane.doe@gmail.com", "Jane Doe", "fake-password")
	invitation, _ = domainmodel.NewProjectInvite(project, user)
	req.ActorID = project.Members[0].ID
	projectInviteMockRepository.
		EXPECT().
		FindById(req.ProjectInviteID).
		Return(invitation, nil)
	userMockRepository.
		EXPECT().
		FindById(project.Members[0].ID).
		Return(project.Members[0], nil)
	projectInviteMockRepository.
		EXPECT().
		DeleteById(req.ProjectInviteID).
		Return(nil)

	err = useCase.Execute(req)
	assert.Nil(t, err)
}
