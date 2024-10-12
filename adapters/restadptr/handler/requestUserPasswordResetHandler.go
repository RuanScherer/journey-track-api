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

type RequestUserPasswordResetHandler struct {
	useCase usecase.RequestUserPasswordResetUseCase
}

func NewRequestUserPasswordResetHandler() *RequestUserPasswordResetHandler {
	userRepository := repository.NewUserPostgresRepository(postgresadptr.GetConnection())
	kafkaProducer := kafkaadptr.NewProducerFactory()
	useCase := *usecase.NewRequestUserPasswordResetUseCase(userRepository, kafkaProducer)
	return &RequestUserPasswordResetHandler{useCase: useCase}
}

func (handler *RequestUserPasswordResetHandler) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.RequestPasswordResetRequest{}
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	err = validator.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	err = handler.useCase.Execute(req)
	return err
}
