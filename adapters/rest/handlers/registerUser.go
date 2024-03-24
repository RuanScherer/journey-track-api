package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgres"
	"github.com/RuanScherer/journey-track-api/adapters/postgres/repositories"
	"github.com/RuanScherer/journey-track-api/adapters/email"
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	"github.com/RuanScherer/journey-track-api/adapters/rest/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type RegisterUserHandler struct {
	useCase usecase.RegisterUserUseCase
}

func NewRegisterUserHandler() *RegisterUserHandler {
	userRepository := repositories.NewUserPostgresRepository(postgres.GetConnection())
	emailService := email.NewSmtpEmailService()
	useCase := *usecase.NewRegisterUserUseCase(userRepository, emailService)
	return &RegisterUserHandler{useCase: useCase}
}

func (handler *RegisterUserHandler) Handle(ctx *fiber.Ctx) error {
	registerUserRequest := &appmodel.RegisterUserRequest{}
	err := ctx.BodyParser(registerUserRequest)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	err = utils.ValidateRequestBody(registerUserRequest)
	if err != nil {
		return err
	}

	registerUserResponse, appErr := handler.useCase.Execute(registerUserRequest)
	if appErr != nil {
		return appErr
	}

	ctx.Status(fiber.StatusCreated)
	return ctx.JSON(registerUserResponse)
}
