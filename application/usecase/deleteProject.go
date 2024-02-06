package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/domain/repository"
	"gorm.io/gorm"
)

type DeleteProjectUseCase struct {
	projectRepository repository.ProjectRepository
}

func NewDeleteProjectUseCase(projectRepository repository.ProjectRepository) *DeleteProjectUseCase {
	return &DeleteProjectUseCase{projectRepository}
}

func (useCase *DeleteProjectUseCase) Execute(req *appmodel.DeleteProjectRequest) error {
	project, err := useCase.projectRepository.FindById(req.ProjectID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appmodel.NewAppError("project_not_found", "project not found", appmodel.ErrorTypeValidation)
		}
		return appmodel.NewAppError("unable_to_find_project", err.Error(), appmodel.ErrorTypeDatabase)
	}

	if project.OwnerID != req.ActorID {
		return appmodel.NewAppError(
			"not_project_owner",
			"only project owner can delete project",
			appmodel.ErrorTypeValidation,
		)
	}

	err = useCase.projectRepository.DeleteById(req.ProjectID)
	if err != nil {
		return appmodel.NewAppError("unable_to_delete_project", err.Error(), appmodel.ErrorTypeDatabase)
	}
	return nil
}
