package usecase

import (
	"errors"
	"github.com/RuanScherer/journey-track-api/application/email"
	"github.com/RuanScherer/journey-track-api/application/factory"
	"github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	domainmodel "github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
	"testing"
)

func TestInviteProjectMembersUseCase_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	projectRepositoryMock := repository.NewMockProjectRepository(ctrl)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	projectInviteRepositoryMock := repository.NewMockProjectInviteRepository(ctrl)
	emailServiceMock := email.NewMockEmailService(ctrl)
	useCase := NewInviteProjectMembersUseCase(
		projectRepositoryMock,
		userRepositoryMock,
		projectInviteRepositoryMock,
		emailServiceMock,
	)

	req := &model.InviteProjectMembersRequest{
		ActorID:   "fake-actor-id",
		ProjectID: "fake-project-id",
		UserIDs:   []string{"fake-user-id-1", "fake-user-id-2"},
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
	assert.Error(t, err, "(database) [unable_to_find_project] unexpected error")

	project, _ := factory.NewProjectWithDefaultOwner("fake project")
	projectRepositoryMock.
		EXPECT().
		FindById(req.ProjectID).
		AnyTimes().
		Return(project, nil)
	userRepositoryMock.
		EXPECT().
		FindById(req.ActorID).
		Return(nil, errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_identify_user] unable to identify the user trying to invite a member")

	actor, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	userRepositoryMock.
		EXPECT().
		FindById(req.ActorID).
		Return(actor, nil)

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [not_project_member] only project members can invite members")

	actor.ID = project.OwnerID
	userRepositoryMock.
		EXPECT().
		FindById(req.ActorID).
		AnyTimes().
		Return(actor, nil)
	userRepositoryMock.
		EXPECT().
		FindById(req.UserIDs[0]).
		Return(nil, gorm.ErrRecordNotFound)

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(validation) [user_not_found] user not found")

	userRepositoryMock.
		EXPECT().
		FindById(req.UserIDs[0]).
		Return(nil, errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_find_user] unexpected error")

	user1, _ := factory.NewVerifiedUser("user1@gmail.com", "User 1", "fake-password")
	user2, _ := factory.NewVerifiedUser("user2@gmail.com", "User 2", "fake-password")
	userRepositoryMock.
		EXPECT().
		FindById(req.UserIDs[0]).
		AnyTimes().
		Return(user1, nil)
	userRepositoryMock.
		EXPECT().
		FindById(req.UserIDs[1]).
		AnyTimes().
		Return(user2, nil)
	projectInviteRepositoryMock.
		EXPECT().
		FindPendingByUserAndProject(gomock.Any(), gomock.Any()).
		Return(nil, errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_check_pending_invites] unexpected error")

	existentInvitation, _ := domainmodel.NewProjectInvite(project, user1)
	projectInviteRepositoryMock.
		EXPECT().
		FindPendingByUserAndProject(user1.ID, project.ID).
		Return(nil, gorm.ErrRecordNotFound)
	projectInviteRepositoryMock.
		EXPECT().
		FindPendingByUserAndProject(user2.ID, project.ID).
		Return(existentInvitation, nil)
	projectInviteRepositoryMock.
		EXPECT().
		BatchCreate(gomock.Any()).
		Return(errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Error(t, err, "(database) [unable_to_save_invites] unexpected error")

	projectInviteRepositoryMock.
		EXPECT().
		FindPendingByUserAndProject(user1.ID, project.ID).
		Return(nil, gorm.ErrRecordNotFound)
	projectInviteRepositoryMock.
		EXPECT().
		FindPendingByUserAndProject(user2.ID, project.ID).
		Return(existentInvitation, nil)
	projectInviteRepositoryMock.
		EXPECT().
		BatchCreate(gomock.Any()).
		Return(nil)
	projectInviteRepositoryMock.
		EXPECT().
		FindById(gomock.Any()).
		AnyTimes().
		Return(nil, errors.New("unexpected error"))

	res, err = useCase.Execute(req)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	invitations := *res
	assert.Len(t, invitations, 2)
	assert.Equal(t, invitations[1].ID, existentInvitation.ID)
	assert.Equal(t, invitations[1].Project.ID, existentInvitation.ProjectID)
	assert.Equal(t, invitations[1].Project.Name, existentInvitation.Project.Name)
	assert.Equal(t, invitations[1].User.ID, existentInvitation.UserID)
	assert.Equal(t, invitations[1].User.Email, *existentInvitation.User.Email)
	assert.Equal(t, invitations[1].User.Name, existentInvitation.User.Name)
	assert.Equal(t, invitations[1].Status, existentInvitation.Status)
}

func TestInviteProjectMembersUseCase_sendProjectInviteEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	projectRepositoryMock := repository.NewMockProjectRepository(ctrl)
	userRepositoryMock := repository.NewMockUserRepository(ctrl)
	projectInviteRepositoryMock := repository.NewMockProjectInviteRepository(ctrl)
	emailServiceMock := email.NewMockEmailService(ctrl)
	useCase := NewInviteProjectMembersUseCase(
		projectRepositoryMock,
		userRepositoryMock,
		projectInviteRepositoryMock,
		emailServiceMock,
	)

	project, _ := factory.NewProjectWithDefaultOwner("fake project")
	user, _ := factory.NewVerifiedUser("john.doe@gmail.com", "John Doe", "fake-password")
	invitation, _ := domainmodel.NewProjectInvite(project, user)

	projectInviteRepositoryMock.
		EXPECT().
		FindById(invitation.ID).
		Return(nil, errors.New("unexpected error"))

	useCase.sendProjectInviteEmail(invitation.ID, user.Name)

	projectInviteRepositoryMock.
		EXPECT().
		FindById(invitation.ID).
		Return(invitation, nil)
	emailServiceMock.
		EXPECT().
		SendEmail(gomock.Any()).
		Return(errors.New("unexpected error"))

	useCase.sendProjectInviteEmail(invitation.ID, user.Name)

	projectInviteRepositoryMock.
		EXPECT().
		FindById(invitation.ID).
		Return(invitation, nil)
	emailServiceMock.
		EXPECT().
		SendEmail(gomock.Any()).
		Return(nil)

	useCase.sendProjectInviteEmail(invitation.ID, user.Name)
}
