package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgres"
	"github.com/RuanScherer/journey-track-api/adapters/postgres/repository"
	"github.com/RuanScherer/journey-track-api/adapters/rest"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type RevokeProjectInviteHandler struct {
	useCase usecase.RevokeProjectInviteUseCase
}

func NewRevokeProjectInviteHandler() *RevokeProjectInviteHandler {
	db := postgres.GetConnection()
	projectInviteRepository := repository.NewProjectInvitePostgresRepository(db)
	userRepository := repository.NewUserPostgresRepository(db)
	useCase := *usecase.NewRevokeProjectInviteUseCase(projectInviteRepository, userRepository)
	return &RevokeProjectInviteHandler{useCase: useCase}
}

func (handler *RevokeProjectInviteHandler) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.RevokeProjectInviteRequest{
		ActorID:         ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ProjectInviteID: ctx.Params("id"),
	}

	err := rest.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	err = handler.useCase.Exceute(req)
	if err != nil {
		return err
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
