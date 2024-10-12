package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr"
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr/repository"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/validator"
	"github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type ShowInvitationByProjectAndTokenHandler struct {
	useCase usecase.ShowInvitationByProjectAndTokenUseCase
}

func NewShowInvitationByProjectAndTokenHandler() *ShowInvitationByProjectAndTokenHandler {
	db := postgresadptr.GetConnection()
	projectInviteRepository := repository.NewProjectInvitePostgresRepository(db)
	useCase := *usecase.NewShowInvitationByProjectAndTokenUseCase(projectInviteRepository)
	return &ShowInvitationByProjectAndTokenHandler{useCase}
}

func (handler *ShowInvitationByProjectAndTokenHandler) Handle(ctx *fiber.Ctx) error {
	req := model.ShowInvitationByProjectAndTokenUseCaseRequest{
		ProjectID: ctx.Params("projectId"),
		Token:     ctx.Params("token"),
	}

	err := validator.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}

	return ctx.JSON(*res)
}
