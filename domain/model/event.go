package model

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type Event struct {
	ID        string     `json:"id" valid:"uuid~[event] Invalid ID"`
	Name      string     `json:"name" valid:"required~[event] Name is required,minstringlength(1)~[event] Name is required,minstringlength(2)~[event] Name should be longer than 2 characters"`
	Timestamp *time.Time `json:"timestamp"`
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
