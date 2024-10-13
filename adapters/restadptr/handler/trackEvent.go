package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr"
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr/repository"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/model"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/validator"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type TrackEventHandler struct {
	useCase usecase.TrackEventUseCase
}

func NewTrackEventHandler() *TrackEventHandler {
	dbConn := postgresadptr.GetConnection()
	projectRepository := repository.NewProjectPostgresRepository(dbConn)
	eventRepository := repository.NewEventPostgresRepository(dbConn)
	useCase := usecase.NewTrackEventUseCase(projectRepository, eventRepository)
	return &TrackEventHandler{*useCase}
}

func (handler *TrackEventHandler) Handle(ctx *fiber.Ctx) error {
	trackEventRequest := &appmodel.TrackEventRequest{}
	if err := ctx.BodyParser(trackEventRequest); err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}
	trackEventRequest.ProjectToken = ctx.Params("projectToken")

	if err := validator.ValidateRequestBody(trackEventRequest); err != nil {
		return err
	}

	if err := handler.useCase.Execute(trackEventRequest); err != nil {
		return err
	}

	ctx.Status(http.StatusNoContent)
	return nil
}
