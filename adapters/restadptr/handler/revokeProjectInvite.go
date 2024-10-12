package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr"
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr/repository"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/validator"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type RevokeProjectInviteHandler struct {
	useCase usecase.RevokeProjectInviteUseCase
}

func NewRevokeProjectInviteHandler() *RevokeProjectInviteHandler {
	db := postgresadptr.GetConnection()
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
