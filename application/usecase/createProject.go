package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/RuanScherer/journey-track-api/domain/repository"
	"gorm.io/gorm"
)

type CreateProjectUseCase struct {
	projectRepository repository.ProjectRepository
	userRepository    repository.UserRepository
}

func NewCreateProjectUseCase(
	projectRepository repository.ProjectRepository,
	userRepository repository.UserRepository,
) *CreateProjectUseCase {
	return &CreateProjectUseCase{
		projectRepository,
		userRepository,
	}
}

func (useCase *CreateProjectUseCase) Execute(req *appmodel.CreateProjectRequest) (*appmodel.CreateProjectResponse, error) {
	ownerUser, err := useCase.userRepository.FindById(req.OwnerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appmodel.NewAppError(
				"project_owner_not_found",
				"project owner not found",
				appmodel.ErrorTypeValidation,
			)
		}
		return nil, appmodel.NewAppError("unable_to_find_project_owner", err.Error(), appmodel.ErrorTypeDatabase)
	}

	project, err := model.NewProject(req.Name, ownerUser)
	if err != nil {
		return nil, appmodel.NewAppError("unable_to_create_project", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.projectRepository.Register(project)
	if err != nil {
		return nil, appmodel.NewAppError("unable_to_save_project", "unable to save project", appmodel.ErrorTypeDatabase)
	}

	return &appmodel.CreateProjectResponse{
		ID:      project.ID,
		Name:    project.Name,
		OwnerID: project.OwnerID,
	}, nil
}
