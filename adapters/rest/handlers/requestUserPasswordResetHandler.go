package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/db"
	"github.com/RuanScherer/journey-track-api/adapters/db/repositories"
	"github.com/RuanScherer/journey-track-api/adapters/email"
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	"github.com/RuanScherer/journey-track-api/adapters/rest/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type RequestUserPasswordResetHandler struct {
	useCase usecase.RequestUserPasswordResetUseCase
}

func NewRequestUserPasswordResetHandler() *RequestUserPasswordResetHandler {
	userRepository := repositories.NewUserDBRepository(db.GetConnection())
	emailService := email.NewSmtpEmailService()
	useCase := *usecase.NewRequestUserPasswordResetUseCase(userRepository, emailService)
	return &RequestUserPasswordResetHandler{useCase: useCase}
}

func (handler *RequestUserPasswordResetHandler) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.RequestPasswordResetRequest{}
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	err = utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	err = handler.useCase.Execute(req)
	return err
}
