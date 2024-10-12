package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/kafkaadptr"
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr"
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr/repository"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/validator"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type InviteProjectMembersHandler struct {
	useCase usecase.InviteProjectMembersUseCase
}

func NewInviteProjectMembersHandler() *InviteProjectMembersHandler {
	db := postgresadptr.GetConnection()
	projectRepository := repository.NewProjectPostgresRepository(db)
	userRepository := repository.NewUserPostgresRepository(db)
	projectInviteRepository := repository.NewProjectInvitePostgresRepository(db)
	producerFactory := kafkaadptr.NewProducerFactory()
	useCase := *usecase.NewInviteProjectMembersUseCase(
		projectRepository,
		userRepository,
		projectInviteRepository,
		producerFactory,
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

	err = validator.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	invite, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(invite)
}
