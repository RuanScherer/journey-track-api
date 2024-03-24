package repository

import (
	"github.com/RuanScherer/journey-track-api/domain/model"
	"gorm.io/gorm"
)

type EventPostgresRepository struct {
	DB *gorm.DB
}

func NewEventPostgresRepository(db *gorm.DB) *EventPostgresRepository {
	return &EventPostgresRepository{DB: db}
}

func (repository *EventPostgresRepository) Register(event *model.Event) error {
	return repository.DB.Create(event).Error
}
