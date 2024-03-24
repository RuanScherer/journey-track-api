package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgres"
	"github.com/RuanScherer/journey-track-api/adapters/postgres/repository"
	"github.com/RuanScherer/journey-track-api/adapters/rest"
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type ResetUserPassword struct {
	useCase usecase.ResetUserPasswordUseCase
}

func NewResetUserPassword() *ResetUserPassword {
	userRepository := repository.NewUserPostgresRepository(postgres.GetConnection())
	useCase := *usecase.NewResetUserPasswordUseCase(userRepository)
	return &ResetUserPassword{useCase: useCase}
}

func (handler *ResetUserPassword) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.PasswordResetRequest{}
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	req.UserID = ctx.Params("id")
	req.PasswordResetToken = ctx.Params("token")

	err = rest.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	err = handler.useCase.Execute(req)
	return err
}
