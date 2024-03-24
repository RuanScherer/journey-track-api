package usecase

import (
	"github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"gorm.io/gorm"
)

type ShowInvitationByProjectAndTokenUseCase struct {
	projectInviteRepository repository.ProjectInviteRepository
}

func NewShowInvitationByProjectAndTokenUseCase(
	projectInviteRepository repository.ProjectInviteRepository,
) *ShowInvitationByProjectAndTokenUseCase {
	return &ShowInvitationByProjectAndTokenUseCase{
		projectInviteRepository,
	}
}

func (useCase *ShowInvitationByProjectAndTokenUseCase) Execute(
	req model.ShowInvitationByProjectAndTokenUseCaseRequest,
) (*model.ShowInvitationByProjectAndTokenUseCaseResponse, error) {
	invitation, err := useCase.projectInviteRepository.FindByProjectAndToken(req.ProjectID, req.Token)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.NewAppError("invitation_not_found", "invitation not found", model.ErrorTypeValidation)
		}
		return nil, model.NewAppError("unable_to_find_invitation", err.Error(), model.ErrorTypeDatabase)
	}

	return &model.ShowInvitationByProjectAndTokenUseCaseResponse{
		ID: invitation.ID,
		Project: &model.InviteProject{
			ID:   invitation.Project.ID,
			Name: invitation.Project.Name,
		},
		User: &model.InviteUser{
			ID:    invitation.User.ID,
			Email: *invitation.User.Email,
			Name:  invitation.User.Name,
		},
		Status: invitation.Status,
	}, nil
}
