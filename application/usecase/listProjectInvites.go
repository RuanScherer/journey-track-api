package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/RuanScherer/journey-track-api/domain/repository"
	"gorm.io/gorm"
)

type ListProjectInvitesUseCase struct {
	projectInviteRepository repository.ProjectInviteRepository
	projectRepository       repository.ProjectRepository
}

func NewListProjectInvitesUseCase(
	projectInviteRepository repository.ProjectInviteRepository,
	projectRepository repository.ProjectRepository,
) *ListProjectInvitesUseCase {
	return &ListProjectInvitesUseCase{
		projectInviteRepository,
		projectRepository,
	}
}

func (useCase *ListProjectInvitesUseCase) Execute(
	req *appmodel.ListProjectInvitesRequest,
) (*appmodel.ListProjectInvitesResponse, error) {
	isMember, err := useCase.projectRepository.HasMember(req.ProjectID, req.ActorID)
	if err != nil {
		return nil, appmodel.NewAppError("unable_to_check_membership", err.Error(), appmodel.ErrorTypeDatabase)
	}

	if !isMember {
		return nil, appmodel.NewAppError(
			"not_project_member",
			"only project members can see the project invites",
			appmodel.ErrorTypeValidation,
		)
	}

	status := req.Status
	if status == "" {
		status = model.ProjectInviteStatusPending
	}

	invites, err := useCase.projectInviteRepository.ListByProjectAndStatus(req.ProjectID, req.Status)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &appmodel.ListProjectInvitesResponse{}, nil
		}
		return nil, appmodel.NewAppError(
			"unable_to_find_project_invites",
			"unable to find the project invites",
			appmodel.ErrorTypeDatabase,
		)
	}

	invitesResponse := []*appmodel.ProjectInvite{}
	for _, invite := range invites {
		invitesResponse = append(invitesResponse, &appmodel.ProjectInvite{
			ID: invite.ID,
			Project: &appmodel.InviteProject{
				ID:   invite.Project.ID,
				Name: invite.Project.Name,
			},
			User: &appmodel.InviteUser{
				ID:    invite.User.ID,
				Name:  invite.User.Name,
				Email: *invite.User.Email,
			},
			Status: invite.Status,
		})
	}
	return &invitesResponse, nil
}
