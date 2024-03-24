package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgres"
	"github.com/RuanScherer/journey-track-api/adapters/postgres/repository"
	"github.com/RuanScherer/journey-track-api/adapters/rest"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type GetProjectStatsHandler struct {
	useCase usecase.GetProjectStatsUseCase
}

func NewGetProjectStatsHandler() *GetProjectStatsHandler {
	db := postgres.GetConnection()
	projectRepository := repository.NewProjectPostgresRepository(db)
	userRepository := repository.NewUserPostgresRepository(db)
	useCase := *usecase.NewGetProjectStatsUseCase(projectRepository, userRepository)
	return &GetProjectStatsHandler{useCase: useCase}
}

func (handler *GetProjectStatsHandler) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.GetProjectStatsRequest{
		ActorID:   ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ProjectID: ctx.Params("id"),
	}

	err := rest.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
