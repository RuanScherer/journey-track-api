package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/kafkaadptr"
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr"
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr/repository"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/model"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/validator"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type RegisterUserHandler struct {
	useCase usecase.RegisterUserUseCase
}

func NewRegisterUserHandler() *RegisterUserHandler {
	userRepository := repository.NewUserPostgresRepository(postgresadptr.GetConnection())
	producerFactory := kafkaadptr.NewProducerFactory()
	useCase := *usecase.NewRegisterUserUseCase(userRepository, producerFactory)
	return &RegisterUserHandler{useCase: useCase}
}

func (handler *RegisterUserHandler) Handle(ctx *fiber.Ctx) error {
	registerUserRequest := &appmodel.RegisterUserRequest{}
	err := ctx.BodyParser(registerUserRequest)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	err = validator.ValidateRequestBody(registerUserRequest)
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
