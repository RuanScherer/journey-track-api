package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/rest/validator"
	"time"

	"github.com/RuanScherer/journey-track-api/adapters/postgres"
	"github.com/RuanScherer/journey-track-api/adapters/postgres/repository"
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	"github.com/RuanScherer/journey-track-api/application/jwt"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type SignInHandler struct {
	useCase usecase.SignInUseCase
}

func NewSignInHandler() *SignInHandler {
	userRepository := repository.NewUserPostgresRepository(postgres.GetConnection())
	jwtManager := jwt.NewDefaultManager()
	useCase := *usecase.NewSignInUseCase(userRepository, jwtManager)
	return &SignInHandler{useCase: useCase}
}

func (handler *SignInHandler) Handle(ctx *fiber.Ctx) error {
	signInRequest := &appmodel.SignInRequest{}
	err := ctx.BodyParser(signInRequest)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	err = validator.ValidateRequestBody(signInRequest)
	if err != nil {
		return err
	}

	signInResponse, appErr := handler.useCase.Execute(signInRequest)
	if appErr != nil {
		return appErr
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    signInResponse.AccessToken,
		HTTPOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(jwt.ExpirationTime),
	})
	return ctx.JSON(signInResponse)
}
