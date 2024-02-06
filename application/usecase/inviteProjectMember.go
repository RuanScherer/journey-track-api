package usecase

import (
	"errors"
	"fmt"
	"log"

	"github.com/RuanScherer/journey-track-api/adapters/email"
	"github.com/RuanScherer/journey-track-api/adapters/email/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/RuanScherer/journey-track-api/domain/repository"
	"github.com/matcornic/hermes/v2"
	"gorm.io/gorm"
)

type InviteProjectMemberUseCase struct {
	projectRepository       repository.ProjectRepository
	userRepository          repository.UserRepository
	projectInviteRepository repository.ProjectInviteRepository
	emailService            email.EmailService
}

func NewInviteProjectMemberUseCase(
	projectRepository repository.ProjectRepository,
	userRepository repository.UserRepository,
	projectInviteRepository repository.ProjectInviteRepository,
	emailService email.EmailService,
) *InviteProjectMemberUseCase {
	return &InviteProjectMemberUseCase{
		projectRepository,
		userRepository,
		projectInviteRepository,
		emailService,
	}
}

func (useCase *InviteProjectMemberUseCase) Execute(
	req *appmodel.InviteProjectMemberRequest,
) (*appmodel.InviteProjectMemberResponse, error) {
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
			"unable to identify the user trying to invite a member",
			appmodel.ErrorTypeDatabase,
		)
	}

	isMember := project.HasMember(actor)
	if !isMember {
		return nil, appmodel.NewAppError(
			"not_project_member",
			"only project members can invite members",
			appmodel.ErrorTypeValidation,
		)
	}

	user, err := useCase.userRepository.FindById(req.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appmodel.NewAppError("user_not_found", "user not found", appmodel.ErrorTypeValidation)
		}
		return nil, appmodel.NewAppError("unable_to_find_user", err.Error(), appmodel.ErrorTypeDatabase)
	}

	existentInvite, err := useCase.projectInviteRepository.FindPendingByUserAndProject(user.ID, project.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appmodel.NewAppError(
			"unable_to_check_pending_invites",
			err.Error(),
			appmodel.ErrorTypeDatabase,
		)
	}

	if existentInvite != nil {
		return nil, appmodel.NewAppError(
			"user_already_invited",
			"user already invited to the project",
			appmodel.ErrorTypeValidation,
		)
	}

	projectInvite, err := model.NewProjectInvite(project, user)
	if err != nil {
		return nil, appmodel.NewAppError("unable_to_invite_user", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.projectInviteRepository.Create(projectInvite)
	if err != nil {
		return nil, appmodel.NewAppError("unable_to_save_invite", err.Error(), appmodel.ErrorTypeDatabase)
	}

	go useCase.sendProjectInviteEmail(projectInvite.ID, actor.Name)
	return &appmodel.InviteProjectMemberResponse{
		ID: projectInvite.ID,
		Project: &appmodel.InviteProject{
			ID:   projectInvite.Project.ID,
			Name: projectInvite.Project.Name,
		},
		User: &appmodel.InviteUser{
			ID:    user.ID,
			Email: *user.Email,
			Name:  user.Name,
		},
		Status: projectInvite.Status,
	}, nil
}

func (useCase *InviteProjectMemberUseCase) sendProjectInviteEmail(inviteId string, issuerName string) {
	invite, err := useCase.projectInviteRepository.FindById(inviteId)
	if err != nil {
		log.Print(err)
		return
	}

	emailConfig := hermes.Email{
		Body: hermes.Body{
			Name:  invite.User.Name,
			Title: "You have been invited to a project",
			Intros: []string{
				fmt.Sprintf("%s has invited you to join the project %s.", issuerName, invite.Project.Name),
				"Join the project to start collaborating with the team.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to answer the invite.",
					Button: hermes.Button{
						Color: "#f25d9c",
						Text:  "Answer invite",
						Link:  "#", // TODO: Add answer invite link when frontend is ready
					},
				},
			},
			Signature: "Regards",
		},
	}
	body, err := utils.GenerateEmailHtml(emailConfig)
	if err != nil {
		log.Print(err)
		return
	}

	err = useCase.emailService.SendEmail(email.EmailSendingConfig{
		To:      *invite.User.Email,
		Subject: "Trackr | You have been invited to a project",
		Body:    body,
	})
	if err != nil {
		log.Print(err)
		return
	}
}
