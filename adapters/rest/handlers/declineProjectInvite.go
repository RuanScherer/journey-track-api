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

type DeclineProjectInviteHandler struct {
	useCase usecase.AcceptProjectInviteUseCase
}

func NewDeclineProjectInviteHandler() *DeclineProjectInviteHandler {
	db := db.GetConnection()
	projectInviteRepository := repository.NewProjectInviteDBRepository(db)
	projectRepository := repository.NewProjectDBRepository(db)
	useCase := *usecase.NewAcceptProjectInviteUseCase(projectInviteRepository, projectRepository)
	return &DeclineProjectInviteHandler{useCase: useCase}
}

func (handler *DeclineProjectInviteHandler) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.AnswerProjectInviteRequest{}
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}
	req.ProjectID = ctx.Params("projectId")

	err = utils.ValidateRequestBody(req)
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