package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/db"
	"github.com/RuanScherer/journey-track-api/adapters/db/repository"
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	"github.com/RuanScherer/journey-track-api/adapters/rest/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type CreateProjectHandler struct {
	useCase usecase.CreateProjectUseCase
}

func NewCreateProjectHandler() *CreateProjectHandler {
	db := db.GetConnection()
	projectRepository := repository.NewProjectDBRepository(db)
	userRepository := repository.NewUserDBRepository(db)
	useCase := *usecase.NewCreateProjectUseCase(projectRepository, userRepository)
	return &CreateProjectHandler{useCase: useCase}
}

func (handler *CreateProjectHandler) Handle(ctx *fiber.Ctx) error {
	req := new(appmodel.CreateProjectRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	err = utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}
