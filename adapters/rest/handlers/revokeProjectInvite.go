package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/db"
	"github.com/RuanScherer/journey-track-api/adapters/db/repositories"
	"github.com/RuanScherer/journey-track-api/adapters/rest/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type RevokeProjectInviteHandler struct {
	useCase usecase.RevokeProjectInviteUseCase
}

func NewRevokeProjectInviteHandler() *RevokeProjectInviteHandler {
	db := db.GetConnection()
	projectInviteRepository := repositories.NewProjectInviteDBRepository(db)
	userRepository := repositories.NewUserDBRepository(db)
	useCase := *usecase.NewRevokeProjectInviteUseCase(projectInviteRepository, userRepository)
	return &RevokeProjectInviteHandler{useCase: useCase}
}

func (handler *RevokeProjectInviteHandler) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.RevokeProjectInviteRequest{
		ActorID:         ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ProjectInviteID: ctx.Params("id"),
	}

	err := utils.ValidateRequestBody(req)
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
