package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgres"
	"github.com/RuanScherer/journey-track-api/adapters/postgres/repository"
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	"github.com/RuanScherer/journey-track-api/adapters/rest/validator"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type EditProjectHandler struct {
	useCase usecase.EditProjectUseCase
}

func NewEditProjectHandler() *EditProjectHandler {
	projectRepository := repository.NewProjectPostgresRepository(postgres.GetConnection())
	useCase := *usecase.NewEditProjectUseCase(projectRepository)
	return &EditProjectHandler{useCase: useCase}
}

func (handler *EditProjectHandler) Handle(ctx *fiber.Ctx) error {
	req := new(appmodel.EditProjectRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	req.ActorID = ctx.Locals("sessionUser").(appmodel.AuthUser).ID
	req.ProjectID = ctx.Params("id")

	err = validator.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
