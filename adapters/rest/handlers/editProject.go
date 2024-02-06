package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/db"
	"github.com/RuanScherer/journey-track-api/adapters/db/repository"
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	"github.com/RuanScherer/journey-track-api/adapters/rest/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type EditProjectHandler struct {
	useCase usecase.EditProjectUseCase
}

func NewEditProjectHandler() *EditProjectHandler {
	projectRepository := repository.NewProjectDBRepository(db.GetConnection())
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

	err = utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
