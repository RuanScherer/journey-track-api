package repository

import "github.com/RuanScherer/journey-track-api/domain/model"

type EventRepository interface {
	Register(*model.Event) error
}
