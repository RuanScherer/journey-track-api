package repository

import (
	"github.com/RuanScherer/journey-track-api/domain/model"
	"gorm.io/gorm"
)

type EventDBRepository struct {
	DB *gorm.DB
}

func NewEventDBRepository(db *gorm.DB) *EventDBRepository {
	return &EventDBRepository{DB: db}
}

func (repository *EventDBRepository) Register(event *model.Event) error {
	return repository.DB.Create(event).Error
}
