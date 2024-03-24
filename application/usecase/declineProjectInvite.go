package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"gorm.io/gorm"
)

type DeclineProjectInviteUseCase struct {
	projectInviteRepository repository.ProjectInviteRepository
}

func NewDeclineProjectInviteUseCase(
	projectInviteRepository repository.ProjectInviteRepository,
) *DeclineProjectInviteUseCase {
	return &DeclineProjectInviteUseCase{projectInviteRepository}
}

func (useCase *DeclineProjectInviteUseCase) Execute(req *appmodel.AnswerProjectInviteRequest) error {
	projectInvite, err := useCase.projectInviteRepository.FindByProjectAndToken(req.ProjectID, req.InviteToken)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appmodel.NewAppError("project_invite_not_found", "project invite not found", appmodel.ErrorTypeValidation)
		}
		return appmodel.NewAppError("unable_to_find_project_invite", err.Error(), appmodel.ErrorTypeDatabase)
	}

	err = projectInvite.Decline(req.InviteToken)
	if err != nil {
		return appmodel.NewAppError("unable_to_decline_project_invite", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.projectInviteRepository.Save(projectInvite)
	if err != nil {
		return appmodel.NewAppError("unable_to_save_project_invite_answer", err.Error(), appmodel.ErrorTypeDatabase)
	}

	return nil
}
