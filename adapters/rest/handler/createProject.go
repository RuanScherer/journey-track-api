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

type CreateProjectHandler struct {
	useCase usecase.CreateProjectUseCase
}

func NewCreateProjectHandler() *CreateProjectHandler {
	db := postgres.GetConnection()
	projectRepository := repository.NewProjectPostgresRepository(db)
	userRepository := repository.NewUserPostgresRepository(db)
	useCase := *usecase.NewCreateProjectUseCase(projectRepository, userRepository)
	return &CreateProjectHandler{useCase: useCase}
}

func (handler *CreateProjectHandler) Handle(ctx *fiber.Ctx) error {
	req := new(appmodel.CreateProjectRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	err = rest.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}
