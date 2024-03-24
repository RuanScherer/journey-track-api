package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/db"
	"github.com/RuanScherer/journey-track-api/adapters/db/repositories"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type ShowUserHandler struct {
	useCase usecase.ShowUserUseCase
}

func NewShowUserHandler() *ShowUserHandler {
	userRepository := repositories.NewUserDBRepository(db.GetConnection())
	useCase := *usecase.NewShowUserUseCase(userRepository)
	return &ShowUserHandler{useCase: useCase}
}

func (handler *ShowUserHandler) Handle(ctx *fiber.Ctx) error {
	userID := ctx.Locals("sessionUser").(appmodel.AuthUser).ID

	response, err := handler.useCase.Execute(userID)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}
