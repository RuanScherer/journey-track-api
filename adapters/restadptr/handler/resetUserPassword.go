package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr"
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr/repository"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/model"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/validator"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type ResetUserPassword struct {
	useCase usecase.ResetUserPasswordUseCase
}

func NewResetUserPassword() *ResetUserPassword {
	userRepository := repository.NewUserPostgresRepository(postgresadptr.GetConnection())
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

	err = validator.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	err = handler.useCase.Execute(req)
	return err
}
