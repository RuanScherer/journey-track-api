package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"

	"github.com/RuanScherer/journey-track-api/adapters/emailtemplateadptr"
	"github.com/RuanScherer/journey-track-api/application/kafka"
	"github.com/RuanScherer/journey-track-api/application/repository"

	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/config"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/matcornic/hermes/v2"
	"gorm.io/gorm"
)

type InviteProjectMembersUseCase struct {
	projectRepository       repository.ProjectRepository
	userRepository          repository.UserRepository
	projectInviteRepository repository.ProjectInviteRepository
	producerFactory         kafka.ProducerFactory
}

func NewInviteProjectMembersUseCase(
	projectRepository repository.ProjectRepository,
	userRepository repository.UserRepository,
	projectInviteRepository repository.ProjectInviteRepository,
	producerFactory kafka.ProducerFactory,
) *InviteProjectMembersUseCase {
	return &InviteProjectMembersUseCase{
		projectRepository,
		userRepository,
		projectInviteRepository,
		producerFactory,
	}
}

func (useCase *InviteProjectMembersUseCase) Execute(
	req *appmodel.InviteProjectMembersRequest,
) (*appmodel.InviteProjectMembersResponse, error) {
	project, err := useCase.projectRepository.FindById(req.ProjectID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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

	invites, e := useCase.generateInvites(req.UserIDs, project)
	if e != nil {
		return nil, e
	}

	// remove invites that already exist from the batch
	invitesToCreate := make([]*model.ProjectInvite, 0)
	for _, invite := range invites {
		if invite.CreatedAt.IsZero() {
			invitesToCreate = append(invitesToCreate, invite)
		}
	}
	err = useCase.projectInviteRepository.BatchCreate(invitesToCreate)
	if err != nil {
		return nil, appmodel.NewAppError("unable_to_save_invites", err.Error(), appmodel.ErrorTypeDatabase)
	}

	useCase.sendProjectInviteEmails(invites, actor)

	var response appmodel.InviteProjectMembersResponse
	for _, invite := range invites {
		response = append(response, &appmodel.ProjectInvite{
			ID: invite.ID,
			Project: &appmodel.InviteProject{
				ID:   invite.Project.ID,
				Name: invite.Project.Name,
			},
			User: &appmodel.InviteUser{
				ID:    invite.User.ID,
				Email: *invite.User.Email,
				Name:  invite.User.Name,
			},
			Status: invite.Status,
		})
	}
	return &response, nil
}

func (useCase *InviteProjectMembersUseCase) generateInvites(
	userIDs []string,
	project *model.Project,
) ([]*model.ProjectInvite, *appmodel.AppError) {
	invites := make([]*model.ProjectInvite, 0)
	for _, userID := range userIDs {
		invite, err := useCase.generateInvite(userID, project)
		if err != nil {
			return make([]*model.ProjectInvite, 0), err
		}
		invites = append(invites, invite)
	}
	return invites, nil
}

func (useCase *InviteProjectMembersUseCase) generateInvite(
	userID string, project *model.Project,
) (*model.ProjectInvite, *appmodel.AppError) {
	user, err := useCase.userRepository.FindById(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
		return existentInvite, nil
	}

	projectInvite, err := model.NewProjectInvite(project, user)
	if err != nil {
		return nil, appmodel.NewAppError("unable_to_invite_user", err.Error(), appmodel.ErrorTypeValidation)
	}
	return projectInvite, nil
}

func (useCase *InviteProjectMembersUseCase) sendProjectInviteEmails(invites []*model.ProjectInvite, actor *model.User) {
	for _, invite := range invites {
		go useCase.queueProjectInviteEmail(invite.ID, actor.Name)
	}
}

func (useCase *InviteProjectMembersUseCase) queueProjectInviteEmail(inviteId string, issuerName string) {
	invite, err := useCase.projectInviteRepository.FindById(inviteId)
	if err != nil {
		slog.Error("Unable to find invite to send email", "error", err)
		return
	}

	appConfig := config.GetAppConfig()
	answerInviteLink := fmt.Sprintf(
		"%s/answer-invitation?projectId=%s&token=%s",
		appConfig.FrontendUrl,
		invite.ProjectID,
		*invite.Token,
	)

	emailTemplate := hermes.Email{
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
						Link:  answerInviteLink,
					},
				},
			},
			Signature: "Regards",
		},
	}
	content, err := emailtemplateadptr.GenerateEmailHtml(emailTemplate)
	if err != nil {
		log.Print(err)
		return
	}

	producer, err := useCase.producerFactory.NewProducer(map[string]any{
		"bootstrap.servers": appConfig.KafkaBootstrapServers,
		"retries":           3,
		"retry.backoff.ms":  1000,
	})
	if err != nil {
		slog.Error("Error setting kafka producer config", "error", err)
		return
	}

	payload, err := json.Marshal(kafka.EmailSendindRequestedPayload{
		To:      *invite.User.Email,
		Subject: "Trackr | You have been invited to a project",
		Content: content,
	})
	if err != nil {
		slog.Error("Error marshalling email sending payload", "error", err)
		return
	}
	message := kafka.Message{Value: payload}
	err = producer.Produce("email-sending-requested", message)
	if err != nil {
		slog.Error("Error producing kafka message to send email", "error", err)
		return
	}
}
