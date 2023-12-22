package usecase

import (
	"errors"

	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"gorm.io/gorm"
)

type ProjectUseCase struct {
	projectRepository       model.ProjectRepository
	userRepository          model.UserRepository
	projectInviteRepository model.ProjectInviteRepository
	eventRepository         model.EventRepository
}

func NewProjectUseCase(
	projectRepository model.ProjectRepository,
	userRepository model.UserRepository,
	projectInviteRepository model.ProjectInviteRepository,
	eventRepository model.EventRepository,
) *ProjectUseCase {
	return &ProjectUseCase{
		projectRepository:       projectRepository,
		userRepository:          userRepository,
		projectInviteRepository: projectInviteRepository,
		eventRepository:         eventRepository,
	}
}

func (useCase *ProjectUseCase) CreateProject(req *appmodel.CreateProjectRequest) (*appmodel.CreateProjectResponse, error) {
	ownerUser, err := useCase.userRepository.FindById(req.OwnerID)
	if err != nil {
		return nil, err
	}

	project, err := model.NewProject(req.Name, ownerUser)
	if err != nil {
		return nil, err
	}

	err = useCase.projectRepository.Register(project)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	err = project.ChangeName(req.Name)
	if err != nil {
		return nil, err
	}

	err = useCase.projectRepository.Save(project)
	if err != nil {
		return nil, err
	}

	return &appmodel.EditProjectResponse{
		ID:      project.ID,
		Name:    project.Name,
		OwnerID: project.OwnerID,
	}, nil
}

func (useCase *ProjectUseCase) ShowProject(projectId string) (*appmodel.ShowProjectResponse, error) {
	project, err := useCase.projectRepository.FindById(projectId)
	if err != nil {
		return nil, err
	}

	var members []*appmodel.ShowUserResponse
	for _, m := range project.Members {
		member := &appmodel.ShowUserResponse{
			ID:         m.ID,
			Email:      m.Email,
			Name:       m.Name,
			IsVerified: m.IsVerified,
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
		return nil, err
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

func (useCase *ProjectUseCase) DeleteProject(projectId string) error {
	err := useCase.projectRepository.DeleteById(projectId)
	return err
}

func (useCase *ProjectUseCase) InviteMember(req *appmodel.InviteProjectMemberRequest) (*appmodel.InviteProjectMemberResponse, error) {
	project, err := useCase.projectRepository.FindById(req.ProjectID)
	if err != nil {
		return nil, err
	}

	user, err := useCase.userRepository.FindById(req.UserID)
	if err != nil {
		return nil, err
	}

	existentInvite, err := useCase.projectInviteRepository.FindPendingByUserAndProject(user.ID, project.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if existentInvite != nil {
		return nil, errors.New("user already invited to the project")
	}

	projectInvite, err := model.NewProjectInvite(project, user)
	if err != nil {
		return nil, err
	}

	err = useCase.projectInviteRepository.Create(projectInvite)
	if err != nil {
		return nil, err
	}

	return &appmodel.InviteProjectMemberResponse{
		ID: projectInvite.ID,
		Project: &appmodel.InviteProject{
			ID:   projectInvite.Project.ID,
			Name: projectInvite.Project.Name,
		},
		User: &appmodel.InviteUser{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
		Status: projectInvite.Status,
	}, nil
}

func (useCase *ProjectUseCase) AcceptInvite(req *appmodel.AnswerProjectInviteRequest) error {
	projectInvite, err := useCase.projectInviteRepository.FindByProjectAndToken(req.ProjectID, req.InviteToken)
	if err != nil {
		return err
	}

	err = projectInvite.Accept(req.InviteToken)
	if err != nil {
		return err
	}

	err = useCase.projectInviteRepository.Save(projectInvite)
	if err != nil {
		return err
	}

	project, err := useCase.projectRepository.FindById(projectInvite.Project.ID)
	if err != nil {
		return err
	}

	err = project.AddMember(projectInvite.User)
	if err != nil {
		return err
	}

	err = useCase.projectRepository.Save(project)
	if err != nil {
		return err
	}

	return nil
}

func (useCase *ProjectUseCase) DeclineInvite(req *appmodel.AnswerProjectInviteRequest) error {
	projectInvite, err := useCase.projectInviteRepository.FindByProjectAndToken(req.ProjectID, req.InviteToken)
	if err != nil {
		return err
	}

	err = projectInvite.Decline(req.InviteToken)
	if err != nil {
		return err
	}

	err = useCase.projectInviteRepository.Save(projectInvite)
	if err != nil {
		return err
	}

	return nil
}

func (useCase *ProjectUseCase) RevokeInvite(projectInviteId string) error {
	projectInvite, err := useCase.projectInviteRepository.FindById(projectInviteId)
	if err != nil {
		return err
	}

	var sessionUser *model.User // TODO: should get authenticated user instead
	isActorProjectMember := projectInvite.Project.HasMember(sessionUser)
	if !isActorProjectMember {
		return errors.New("only members of the project can revoke invites")
	}

	err = useCase.projectInviteRepository.DeleteById(projectInviteId)
	if err != nil {
		return err
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
