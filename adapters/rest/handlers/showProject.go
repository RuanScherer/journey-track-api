package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgres"
	"github.com/RuanScherer/journey-track-api/adapters/postgres/repositories"
	"github.com/RuanScherer/journey-track-api/adapters/rest/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type ShowProjectHandler struct {
	useCase usecase.ShowProjectUseCase
}

func NewShowProjectHandler() *ShowProjectHandler {
	db := postgres.GetConnection()
	projectRepository := repositories.NewProjectPostgresRepository(db)
	userRepository := repositories.NewUserPostgresRepository(db)
	useCase := *usecase.NewShowProjectUseCase(projectRepository, userRepository)
	return &ShowProjectHandler{useCase: useCase}
}

func (handler *ShowProjectHandler) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.ShowProjectRequest{
		ActorID:   ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ProjectID: ctx.Params("id"),
	}

	err := utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
