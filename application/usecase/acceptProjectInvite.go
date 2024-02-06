package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/domain/repository"
	"gorm.io/gorm"
)

type AcceptProjectInviteUseCase struct {
	projectInviteRepository repository.ProjectInviteRepository
	projectRepository       repository.ProjectRepository
}

func NewAcceptProjectInviteUseCase(
	projectInviteRepository repository.ProjectInviteRepository,
	projectRepository repository.ProjectRepository,
) *AcceptProjectInviteUseCase {
	return &AcceptProjectInviteUseCase{
		projectInviteRepository,
		projectRepository,
	}
}

func (useCase *AcceptProjectInviteUseCase) Execute(req *appmodel.AnswerProjectInviteRequest) error {
	projectInvite, err := useCase.projectInviteRepository.FindByProjectAndToken(req.ProjectID, req.InviteToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appmodel.NewAppError("project_invite_not_found", "project invite not found", appmodel.ErrorTypeValidation)
		}
		return appmodel.NewAppError("unable_to_find_project_invite", err.Error(), appmodel.ErrorTypeDatabase)
	}

	err = projectInvite.Accept(req.InviteToken)
	if err != nil {
		return appmodel.NewAppError("unable_to_accept_project_invite", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.projectInviteRepository.Save(projectInvite)
	if err != nil {
		return appmodel.NewAppError("unable_to_save_project_invite_answer", err.Error(), appmodel.ErrorTypeDatabase)
	}

	project, err := useCase.projectRepository.FindById(projectInvite.Project.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appmodel.NewAppError("project_not_found", "project not found", appmodel.ErrorTypeValidation)
		}
		return appmodel.NewAppError("unable_to_find_project", err.Error(), appmodel.ErrorTypeDatabase)
	}

	err = project.AddMember(projectInvite.User)
	if err != nil {
		return appmodel.NewAppError("unable_to_add_project_member", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.projectRepository.Save(project)
	if err != nil {
		return appmodel.NewAppError("unable_to_save_project_changes", err.Error(), appmodel.ErrorTypeDatabase)
	}

	return nil
}
