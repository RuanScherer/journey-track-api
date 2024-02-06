package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/db"
	"github.com/RuanScherer/journey-track-api/adapters/db/repository"
	"github.com/RuanScherer/journey-track-api/adapters/email"
	"github.com/RuanScherer/journey-track-api/adapters/rest/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type InviteProjectMemberHandler struct {
	useCase usecase.InviteProjectMemberUseCase
}

func NewInviteProjectMemberHandler() *InviteProjectMemberHandler {
	db := db.GetConnection()
	projectRepository := repository.NewProjectDBRepository(db)
	userRepository := repository.NewUserDBRepository(db)
	projectInviteRepository := repository.NewProjectInviteDBRepository(db)
	emailService := email.NewSmtpEmailService()
	useCase := *usecase.NewInviteProjectMemberUseCase(
		projectRepository,
		userRepository,
		projectInviteRepository,
		emailService,
	)
	return &InviteProjectMemberHandler{useCase: useCase}
}

func (handler *InviteProjectMemberHandler) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.InviteProjectMemberRequest{
		ActorID:   ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ProjectID: ctx.Params("projectId"),
		UserID:    ctx.Params("userId"),
	}

	err := utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	invite, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(invite)
}
