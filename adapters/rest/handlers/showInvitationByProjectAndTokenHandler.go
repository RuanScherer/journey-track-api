package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgres"
	"github.com/RuanScherer/journey-track-api/adapters/postgres/repositories"
	"github.com/RuanScherer/journey-track-api/adapters/rest/utils"
	"github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type ShowInvitationByProjectAndTokenHandler struct {
	useCase usecase.ShowInvitationByProjectAndTokenUseCase
}

func NewShowInvitationByProjectAndTokenHandler() *ShowInvitationByProjectAndTokenHandler {
	db := postgres.GetConnection()
	projectInviteRepository := repositories.NewProjectInvitePostgresRepository(db)
	useCase := *usecase.NewShowInvitationByProjectAndTokenUseCase(projectInviteRepository)
	return &ShowInvitationByProjectAndTokenHandler{useCase}
}

func (handler *ShowInvitationByProjectAndTokenHandler) Handle(ctx *fiber.Ctx) error {
	req := model.ShowInvitationByProjectAndTokenUseCaseRequest{
		ProjectID: ctx.Params("projectId"),
		Token:     ctx.Params("token"),
	}

	err := utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}

	return ctx.JSON(*res)
}
