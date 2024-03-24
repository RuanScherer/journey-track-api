package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgres"
	"github.com/RuanScherer/journey-track-api/adapters/postgres/repository"
	"github.com/RuanScherer/journey-track-api/adapters/rest"
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	"github.com/RuanScherer/journey-track-api/adapters/smtpemail"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type RegisterUserHandler struct {
	useCase usecase.RegisterUserUseCase
}

func NewRegisterUserHandler() *RegisterUserHandler {
	userRepository := repository.NewUserPostgresRepository(postgres.GetConnection())
	emailService := smtpemail.NewSmtpEmailService()
	useCase := *usecase.NewRegisterUserUseCase(userRepository, emailService)
	return &RegisterUserHandler{useCase: useCase}
}

func (handler *RegisterUserHandler) Handle(ctx *fiber.Ctx) error {
	registerUserRequest := &appmodel.RegisterUserRequest{}
	err := ctx.BodyParser(registerUserRequest)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	err = rest.ValidateRequestBody(registerUserRequest)
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
