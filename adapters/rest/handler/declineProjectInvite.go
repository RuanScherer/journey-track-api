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

type DeclineProjectInviteHandler struct {
	useCase usecase.DeclineProjectInviteUseCase
}

func NewDeclineProjectInviteHandler() *DeclineProjectInviteHandler {
	db := postgres.GetConnection()
	projectInviteRepository := repository.NewProjectInvitePostgresRepository(db)
	useCase := *usecase.NewDeclineProjectInviteUseCase(projectInviteRepository)
	return &DeclineProjectInviteHandler{useCase: useCase}
}

func (handler *DeclineProjectInviteHandler) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.AnswerProjectInviteRequest{}
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}
	req.ProjectID = ctx.Params("projectId")

	err = validator.ValidateRequestBody(req)
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
