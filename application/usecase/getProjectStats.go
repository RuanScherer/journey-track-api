package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	repository2 "github.com/RuanScherer/journey-track-api/application/repository"
)

type GetProjectStatsUseCase struct {
	projectRepository repository2.ProjectRepository
	userRepository    repository2.UserRepository
}

func NewGetProjectStatsUseCase(
	projectRepository repository2.ProjectRepository,
	userRepository repository2.UserRepository,
) *GetProjectStatsUseCase {
	return &GetProjectStatsUseCase{
		projectRepository,
		userRepository,
	}
}

func (useCase *GetProjectStatsUseCase) Execute(
	req *appmodel.GetProjectStatsRequest,
) (*appmodel.GetProjectStatsResponse, error) {
	isMember, err := useCase.projectRepository.HasMember(req.ProjectID, req.ActorID)
	if err != nil {
		return nil, appmodel.NewAppError("unable_to_check_membership", err.Error(), appmodel.ErrorTypeDatabase)
	}

	if !isMember {
		return nil, appmodel.NewAppError(
			"not_project_member",
			"only project members can see project details",
			appmodel.ErrorTypeValidation,
		)
	}

	stats, err := useCase.projectRepository.FindMembersCountAndEventsCountById(req.ProjectID)
	if err != nil {
		return nil, appmodel.NewAppError("unable_to_load_stats", err.Error(), appmodel.ErrorTypeDatabase)
	}

	return &appmodel.GetProjectStatsResponse{
		InvitesCount: stats.InvitesCount,
		EventsCount:  stats.EventsCount,
	}, nil
}
