package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/domain/repository"
	"gorm.io/gorm"
)

type ShowProjectUseCase struct {
	projectRepository repository.ProjectRepository
	userRepository    repository.UserRepository
}

func NewShowProjectUseCase(
	projectRepository repository.ProjectRepository,
	userRepository repository.UserRepository,
) *ShowProjectUseCase {
	return &ShowProjectUseCase{
		projectRepository,
		userRepository,
	}
}

func (useCase *ShowProjectUseCase) Execute(req *appmodel.ShowProjectRequest) (*appmodel.ShowProjectResponse, error) {
	project, err := useCase.projectRepository.FindById(req.ProjectID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appmodel.NewAppError("project_not_found", "project not found", appmodel.ErrorTypeValidation)
		}
		return nil, appmodel.NewAppError("unable_to_find_project", err.Error(), appmodel.ErrorTypeDatabase)
	}

	actor, err := useCase.userRepository.FindById(req.ActorID)
	if err != nil {
		return nil, appmodel.NewAppError(
			"unable_to_identify_user",
			"unable to identify the user trying to see project details",
			appmodel.ErrorTypeDatabase,
		)
	}

	isMember := project.HasMember(actor)
	if !isMember {
		return nil, appmodel.NewAppError(
			"not_project_member",
			"only project members can see project details",
			appmodel.ErrorTypeValidation,
		)
	}

	var members []*appmodel.ProjectMember
	for _, m := range project.Members {
		member := &appmodel.ProjectMember{
			ID:    m.ID,
			Email: *m.Email,
			Name:  m.Name,
		}
		members = append(members, member)
	}

	return &appmodel.ShowProjectResponse{
		ID:      project.ID,
		Name:    project.Name,
		OwnerID: project.OwnerID,
		IsOwner: project.OwnerID == req.ActorID,
		Members: members,
	}, nil
}
