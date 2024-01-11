package usecase

import (
	"errors"
	"log"

	"github.com/RuanScherer/journey-track-api/adapters/email"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/utils"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"gorm.io/gorm"
)

type ProjectUseCase struct {
	projectRepository       model.ProjectRepository
	userRepository          model.UserRepository
	projectInviteRepository model.ProjectInviteRepository
	eventRepository         model.EventRepository
	emailService            email.EmailService
}

func NewProjectUseCase(
	projectRepository model.ProjectRepository,
	userRepository model.UserRepository,
	projectInviteRepository model.ProjectInviteRepository,
	eventRepository model.EventRepository,
	emailService email.EmailService,
) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepository:       projectRepository,
		userRepository:          userRepository,
		projectInviteRepository: projectInviteRepository,
		eventRepository:         eventRepository,
		emailService:            emailService,
	}
}

func (useCase *ProjectUseCase) CreateProject(req *appmodel.CreateProjectRequest) (*appmodel.CreateProjectResponse, error) {
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

func (useCase *ProjectUseCase) EditProject(req *appmodel.EditProjectRequest) (*appmodel.EditProjectResponse, error) {
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

func (useCase *ProjectUseCase) ShowProject(req *appmodel.ShowProjectRequest) (*appmodel.ShowProjectResponse, error) {
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
		Members: members,
	}, nil
}

func (useCase *ProjectUseCase) ListProjectsByMember(memberId string) (*appmodel.ListProjectByMemberResponse, error) {
	projects, err := useCase.projectRepository.FindByMemberId(memberId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &appmodel.ListProjectByMemberResponse{}, nil
		}
		return nil, appmodel.NewAppError("unable_to_find_projects", err.Error(), appmodel.ErrorTypeDatabase)
	}

	var projectsResponse appmodel.ListProjectByMemberResponse
	for _, p := range projects {
		project := &appmodel.ProjectByMember{
			ID:      p.ID,
			Name:    p.Name,
			OwnerID: p.OwnerID,
		}
		projectsResponse = append(projectsResponse, project)
	}

	return &projectsResponse, nil
}

func (useCase *ProjectUseCase) DeleteProject(req *appmodel.DeleteProjectRequest) error {
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

func (useCase *ProjectUseCase) InviteMember(req *appmodel.InviteProjectMemberRequest) (*appmodel.InviteProjectMemberResponse, error) {
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

func (useCase *ProjectUseCase) sendProjectInviteEmail(inviteId string, issuerName string) {
	invite, err := useCase.projectInviteRepository.FindById(inviteId)
	if err != nil {
		log.Print(err)
		return
	}

	body, err := utils.GetFilledEmailTemplate("project_invite.html", appmodel.ProjectInviteEmailConfig{
		UserName:         invite.User.Name,
		IssuerName:       issuerName,
		ProjectName:      invite.Project.Name,
		AnswerInviteLink: "#", // TODO: Add answer invite link when frontend is ready
	})
	if err != nil {
		log.Print(err)
		return
	}

	err = useCase.emailService.SendEmail(email.EmailSendingConfig{
		To:      *invite.User.Email,
		Subject: "Journey Track | Você foi convidado(a) para um projeto",
		Body:    body,
	})
	if err != nil {
		log.Print(err)
		return
	}
}

func (useCase *ProjectUseCase) AcceptInvite(req *appmodel.AnswerProjectInviteRequest) error {
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

func (useCase *ProjectUseCase) DeclineInvite(req *appmodel.AnswerProjectInviteRequest) error {
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

func (useCase *ProjectUseCase) RevokeInvite(req *appmodel.RevokeProjectInviteRequest) error {
	projectInvite, err := useCase.projectInviteRepository.FindById(req.ProjectInviteID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
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
			"unable to identify the user trying to see project details",
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

/*
TODO: this use case needs to be refactored to implement some security rules like
validating token, requiring sdk authentication, etc.
*/
func (useCase *ProjectUseCase) RegisterEvent(req *appmodel.RegisterEventRequest) error {
	project, err := useCase.projectRepository.FindById(req.ProjectID)
	if err != nil {
		return err
	}

	event, err := model.NewEvent(req.Name, project)
	if err != nil {
		return err
	}

	err = useCase.eventRepository.Register(event)
	return err
}
