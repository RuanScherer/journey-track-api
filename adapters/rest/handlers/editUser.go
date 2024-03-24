package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgres"
	"github.com/RuanScherer/journey-track-api/adapters/postgres/repositories"
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	"github.com/RuanScherer/journey-track-api/adapters/rest/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type EditUserHandler struct {
	useCase usecase.EditUserUseCase
}

func NewEditUserHandler() *EditUserHandler {
	userRepository := repositories.NewUserPostgresRepository(postgres.GetConnection())
	useCase := *usecase.NewEditUserUseCase(userRepository)
	return &EditUserHandler{useCase: useCase}
}

func (handler *EditUserHandler) Handle(ctx *fiber.Ctx) error {
	editUserRequest := &appmodel.EditUserRequest{}
	err := ctx.BodyParser(editUserRequest)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}
	editUserRequest.UserID = ctx.Locals("sessionUser").(appmodel.AuthUser).ID

	err = utils.ValidateRequestBody(editUserRequest)
	if err != nil {
		return err
	}

	response, appErr := handler.useCase.Execute(editUserRequest)
	if appErr != nil {
		return appErr
	}
	return ctx.JSON(response)
}
