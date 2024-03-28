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

func TestDeclineProjectInviteUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	projectInviteRepositoryMock := repository.NewMockProjectInviteRepository(ctrl)
	useCase := NewDeclineProjectInviteUseCase(projectInviteRepositoryMock)

	req := &model.AnswerProjectInviteRequest{
		ProjectID:   "fake-project-id",
		InviteToken: "fake-invite-token",
	}

	projectInviteRepositoryMock.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(nil, gorm.ErrRecordNotFound)

	err := useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [project_invite_not_found] project invite not found")

	projectInviteRepositoryMock.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(nil, errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_find_project_invite] unexpected error")

	project, _ := factory.NewProjectWithDefaultOwner("fake project")
	user, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Dode", "fake-password")
	invitation, _ := domainmodel.NewProjectInvite(project, user)
	projectInviteRepositoryMock.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(invitation, nil)

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_decline_project_invite] invalid token provided to answer invite")

	invitation, _ = domainmodel.NewProjectInvite(project, user)
	req.InviteToken = *invitation.Token
	projectInviteRepositoryMock.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(invitation, nil)
	projectInviteRepositoryMock.
		EXPECT().
		Save(invitation).
		Return(errors.New("unexpected error"))

	err = useCase.Execute(req)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_save_project_invite_answer] unexpected error")

	invitation, _ = domainmodel.NewProjectInvite(project, user)
	req.InviteToken = *invitation.Token
	projectInviteRepositoryMock.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.InviteToken).
		Return(invitation, nil)
	projectInviteRepositoryMock.
		EXPECT().
		Save(invitation).
		Return(nil)

	err = useCase.Execute(req)
	assert.Nil(t, err)
}
