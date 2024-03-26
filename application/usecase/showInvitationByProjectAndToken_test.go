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

func TestShowInvitationByProjectAndTokenUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	projectInviteMockRepository := repository.NewMockProjectInviteRepository(ctrl)
	useCase := NewShowInvitationByProjectAndTokenUseCase(projectInviteMockRepository)

	req := model.ShowInvitationByProjectAndTokenUseCaseRequest{
		ProjectID: "fake-project-id",
		Token:     "fake-token",
	}

	projectInviteMockRepository.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.Token).
		Return(nil, gorm.ErrRecordNotFound)

	res, err := useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [invitation_not_found] invitation not found")

	projectInviteMockRepository.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.Token).
		Return(nil, errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [unable_to_find_invitation] unexpected error")

	project, _ := factory.NewProjectWithDefaultOwner("fake-project-name")
	user, _ := factory.NewVerifiedUser("jane.doe@gmail.com", "Jane Doe", "fake-password")
	invitation, _ := domainmodel.NewProjectInvite(project, user)
	projectInviteMockRepository.
		EXPECT().
		FindByProjectAndToken(req.ProjectID, req.Token).
		Return(invitation, nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, res.ID, invitation.ID)
	assert.Equal(t, res.Project.ID, invitation.Project.ID)
	assert.Equal(t, res.Project.Name, invitation.Project.Name)
	assert.Equal(t, res.User.ID, invitation.User.ID)
	assert.Equal(t, res.User.Email, *invitation.User.Email)
	assert.Equal(t, res.User.Name, invitation.User.Name)
	assert.Equal(t, res.Status, invitation.Status)
}
