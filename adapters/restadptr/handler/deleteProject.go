package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr"
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr/repository"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/validator"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type DeleteProjectHandler struct {
	useCase usecase.DeleteProjectUseCase
}

func NewDeleteProjectHandler() *DeleteProjectHandler {
	projectRepository := repository.NewProjectPostgresRepository(postgresadptr.GetConnection())
	useCase := *usecase.NewDeleteProjectUseCase(projectRepository)
	return &DeleteProjectHandler{useCase: useCase}
}

func (handler *DeleteProjectHandler) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.DeleteProjectRequest{
		ActorID:   ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ProjectID: ctx.Params("id"),
	}

	err := validator.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	err = handler.useCase.Execute(req)
	if err != nil {
		return err
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
