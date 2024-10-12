package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr"
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr/repository"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/validator"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type VerifyUserHandler struct {
	useCase usecase.VerifyUserUseCase
}

func NewVerifyUserHandler() *VerifyUserHandler {
	userRepository := repository.NewUserPostgresRepository(postgresadptr.GetConnection())
	useCase := *usecase.NewVerifyUserUseCase(userRepository)
	return &VerifyUserHandler{useCase: useCase}
}

func (handler *VerifyUserHandler) Handle(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	token := ctx.Params("token")
	verifyUserRequest := &appmodel.VerifyUserRequest{
		UserID:            userId,
		VerificationToken: token,
	}

	err := validator.ValidateRequestBody(verifyUserRequest)
	if err != nil {
		return err
	}

	appErr := handler.useCase.Execute(verifyUserRequest)
	if appErr != nil {
		return appErr
	}

	ctx.Status(fiber.StatusOK)
	return nil
}
