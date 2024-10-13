package usecase

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/repository"
	"github.com/RuanScherer/journey-track-api/domain/model"
)

type TrackEventUseCase struct {
	projectRepository repository.ProjectRepository
	eventRepository   repository.EventRepository
}

func NewTrackEventUseCase(
	projectRepository repository.ProjectRepository,
	eventRepository repository.EventRepository,
) *TrackEventUseCase {
	return &TrackEventUseCase{projectRepository, eventRepository}
}

func (useCase *TrackEventUseCase) Execute(req *appmodel.TrackEventRequest) error {
	project, err := useCase.projectRepository.FindByToken(req.ProjectToken)
	if err != nil {
		return appmodel.NewAppError("project_not_found", "project not found", appmodel.ErrorTypeDatabase)
	}

	event, err := model.NewEvent(req.Name, project)
	if err != nil {
		return appmodel.NewAppError("invalid_data_to_track_event", err.Error(), appmodel.ErrorTypeValidation)
	}

	err = useCase.eventRepository.Register(event)
	if err != nil {
		return appmodel.NewAppError("unable_to_track_event", err.Error(), appmodel.ErrorTypeDatabase)
	}

	return nil
}
