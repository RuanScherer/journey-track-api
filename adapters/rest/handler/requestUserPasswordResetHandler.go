package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgres"
	"github.com/RuanScherer/journey-track-api/adapters/postgres/repository"
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	"github.com/RuanScherer/journey-track-api/adapters/rest/validator"
	"github.com/RuanScherer/journey-track-api/adapters/smtpemail"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type RequestUserPasswordResetHandler struct {
	useCase usecase.RequestUserPasswordResetUseCase
}

func NewRequestUserPasswordResetHandler() *RequestUserPasswordResetHandler {
	userRepository := repository.NewUserPostgresRepository(postgres.GetConnection())
	emailService := smtpemail.NewSmtpEmailService()
	useCase := *usecase.NewRequestUserPasswordResetUseCase(userRepository, emailService)
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
