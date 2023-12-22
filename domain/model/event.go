package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository interface {
	Register(*Event) error
}

type Event struct {
	gorm.Model
	ID        string     `json:"id" gorm:"primaryKey" valid:"uuid~[event] Invalid ID"`
	Name      string     `json:"name" gorm:"type:varchar(255);not null" valid:"required~[event] Name is required,minstringlength(1)~[event] Name is required,minstringlength(2)~[event] Name should be longer than 2 characters"`
	Timestamp *time.Time `json:"timestamp" gorm:"type:timestamp with time zone;not null;default:NOW()" valid:"required~[event] Timestamp is required"`
	ProjectID string     `json:"project_id" gorm:"column:project_id;type:varchar(255);not null" valid:"-"`
	Project   *Project   `json:"project" valid:"-"`
}

func NewEvent(name string, project *Project) (*Event, error) {
	_, err := govalidator.ValidateStruct(project)
	if err != nil {
		return nil, err
	}

	eventTimestamp := time.Now()
	event := &Event{
		ID:        uuid.New().String(),
		Name:      name,
		Timestamp: &eventTimestamp,
		Project:   project,
	}

	_, err = govalidator.ValidateStruct(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}
