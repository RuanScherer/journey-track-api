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

func TestAcceptProjectInviteUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	projectInviteMockRepository := repository.NewMockProjectInviteRepository(ctrl)
	projectMockRepository := repository.NewMockProjectRepository(ctrl)
	useCase := NewAcceptProjectInviteUseCase(projectInviteMockRepository, projectMockRepository)

	req := &model.AnswerProjectInviteRequest{
		ProjectID:   "fake-project-id",
		InviteToken: "fake-invite-token",
	}

	projectInviteMockRepository.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(nil, gorm.ErrRecordNotFound)

	err := useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [project_invite_not_found] project invite not found")

	projectInviteMockRepository.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(nil, errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_find_project_invite] unexpected error")

	project, _ := factory.NewProjectWithDefaultOwner("fake project")
	user, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	invitation, _ := domainmodel.NewProjectInvite(project, user)
	projectInviteMockRepository.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(invitation, nil)

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_accept_project_invite] invalid token provided to answer invite")

	invitation, _ = domainmodel.NewProjectInvite(project, user)
	req.InviteToken = *invitation.Token
	projectInviteMockRepository.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(invitation, nil)
	projectInviteMockRepository.
		EXPECT().
		Save(invitation).
		Return(errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_save_project_invite_answer] unexpected error")

	invitation, _ = domainmodel.NewProjectInvite(project, user)
	req.InviteToken = *invitation.Token
	projectInviteMockRepository.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(invitation, nil)
	projectInviteMockRepository.
		EXPECT().
		Save(invitation).
		Return(nil)
	projectMockRepository.
		EXPECT().
		FindById(project.ID).
		Return(nil, gorm.ErrRecordNotFound)

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [project_not_found] project not found")

	invitation, _ = domainmodel.NewProjectInvite(project, user)
	req.InviteToken = *invitation.Token
	projectInviteMockRepository.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(invitation, nil)
	projectInviteMockRepository.
		EXPECT().
		Save(invitation).
		Return(nil)
	projectMockRepository.
		EXPECT().
		FindById(project.ID).
		Return(nil, errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_find_project] unexpected error")

	invitation, _ = domainmodel.NewProjectInvite(project, user)
	req.InviteToken = *invitation.Token
	invalidProject := *project
	_ = invalidProject.AddMember(user)
	projectInviteMockRepository.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(invitation, nil)
	projectInviteMockRepository.
		EXPECT().
		Save(invitation).
		Return(nil)
	projectMockRepository.
		EXPECT().
		FindById(project.ID).
		Return(&invalidProject, nil)

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_add_project_member] user is already a member of this project")

	invitation, _ = domainmodel.NewProjectInvite(project, user)
	req.InviteToken = *invitation.Token
	projectInviteMockRepository.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(invitation, nil)
	projectInviteMockRepository.
		EXPECT().
		Save(invitation).
		Return(nil)
	projectMockRepository.
		EXPECT().
		FindById(project.ID).
		Return(project, nil)
	projectMockRepository.
		EXPECT().
		Save(project).
		Return(errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_save_project_changes] unexpected error")

	project, _ = factory.NewProjectWithDefaultOwner("fake project")
	invitation, err = domainmodel.NewProjectInvite(project, user)
	req.InviteToken = *invitation.Token
	projectInviteMockRepository.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(invitation, nil)
	projectInviteMockRepository.
		EXPECT().
		Save(invitation).
		Return(nil)
	projectMockRepository.
		EXPECT().
		FindById(project.ID).
		Return(project, nil)
	projectMockRepository.
		EXPECT().
		Save(project).
		Return(nil)

	err = useCase.Execute(req)
	assert.Nil(t, err)
}
