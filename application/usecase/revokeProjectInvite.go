package usecase

import (
	"errors"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	repository2 "github.com/RuanScherer/journey-track-api/application/repository"
	"gorm.io/gorm"
)

type RevokeProjectInviteUseCase struct {
	projectInviteRepository repository2.ProjectInviteRepository
	userRepository          repository2.UserRepository
}

func NewRevokeProjectInviteUseCase(
	projectInviteRepository repository2.ProjectInviteRepository,
	userRepository repository2.UserRepository,
) *RevokeProjectInviteUseCase {
	return &RevokeProjectInviteUseCase{
		projectInviteRepository,
		userRepository,
	}
}

func (useCase *RevokeProjectInviteUseCase) Execute(req *appmodel.RevokeProjectInviteRequest) error {
	projectInvite, err := useCase.projectInviteRepository.FindById(req.ProjectInviteID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return appmodel.NewAppError("project_invite_not_found", "project invite not found", appmodel.ErrorTypeValidation)
		}
		return appmodel.NewAppError("unable_to_find_project_invite", err.Error(), appmodel.ErrorTypeDatabase)
	}

	canRevoke, reason := projectInvite.CanRevoke()
	if !canRevoke {
		return appmodel.NewAppError(
			"unable_to_revoke_project_invite",
			reason,
			appmodel.ErrorTypeValidation,
		)
	}

	actor, err := useCase.userRepository.FindById(req.ActorID)
	if err != nil {
		return appmodel.NewAppError(
			"unable_to_identify_user",
			"unable to identify the user trying to revoke the invite",
			appmodel.ErrorTypeDatabase,
		)
	}

	isActorProjectMember := projectInvite.Project.HasMember(actor)
	if !isActorProjectMember {
		return appmodel.NewAppError(
			"not_project_member",
			"only project members can revoke invites",
			appmodel.ErrorTypeValidation,
		)
	}

	err = useCase.projectInviteRepository.DeleteById(req.ProjectInviteID)
	if err != nil {
		return appmodel.NewAppError("unable_to_delete_project_invite", err.Error(), appmodel.ErrorTypeDatabase)
	}
	return nil
}
