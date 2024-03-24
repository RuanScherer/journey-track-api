package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/db"
	"github.com/RuanScherer/journey-track-api/adapters/db/repositories"
	"github.com/RuanScherer/journey-track-api/adapters/email"
	"github.com/RuanScherer/journey-track-api/adapters/rest/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type InviteProjectMembersHandler struct {
	useCase usecase.InviteProjectMembersUseCase
}

func NewInviteProjectMembersHandler() *InviteProjectMembersHandler {
	db := db.GetConnection()
	projectRepository := repositories.NewProjectDBRepository(db)
	userRepository := repositories.NewUserDBRepository(db)
	projectInviteRepository := repositories.NewProjectInviteDBRepository(db)
	emailService := email.NewSmtpEmailService()
	useCase := *usecase.NewInviteProjectMembersUseCase(
		projectRepository,
		userRepository,
		projectInviteRepository,
		emailService,
	)
	return &InviteProjectMembersHandler{useCase: useCase}
}

func (handler *InviteProjectMembersHandler) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.InviteProjectMembersRequest{}
	err := ctx.BodyParser(req)
	if err != nil {
		return err
	}

	req.ActorID = ctx.Locals("sessionUser").(appmodel.AuthUser).ID
	req.ProjectID = ctx.Params("projectId")

	err = utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	invite, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(invite)
}
