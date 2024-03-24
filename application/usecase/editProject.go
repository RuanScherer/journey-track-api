package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"gorm.io/gorm"
)

type EditProjectUseCase struct {
	projectRepository repository.ProjectRepository
}

func NewEditProjectUseCase(projectRepository repository.ProjectRepository) *EditProjectUseCase {
	return &EditProjectUseCase{projectRepository}
}

func (useCase *EditProjectUseCase) Execute(req *appmodel.EditProjectRequest) (*appmodel.EditProjectResponse, error) {
	project, err := useCase.projectRepository.FindById(req.ProjectID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appmodel.NewAppError("project_not_found", "project not found", appmodel.ErrorTypeValidation)
		}
		return nil, appmodel.NewAppError("unable_to_find_project", err.Error(), appmodel.ErrorTypeDatabase)
	}

	if project.OwnerID != req.ActorID {
		return nil, appmodel.NewAppError(
			"not_project_owner",
			"only project owner can edit project",
			appmodel.ErrorTypeValidation,
		)
	}

	err = project.ChangeName(req.Name)
	if err != nil {
		return nil, appmodel.NewAppError("unable_to_edit_project", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.projectRepository.Save(project)
	if err != nil {
		return nil, appmodel.NewAppError(
			"unable_to_save_project_changes",
			"unable to save project changes",
			appmodel.ErrorTypeDatabase,
		)
	}

	return &appmodel.EditProjectResponse{
		ID:      project.ID,
		Name:    project.Name,
		OwnerID: project.OwnerID,
	}, nil
}
