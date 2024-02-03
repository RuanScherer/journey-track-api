package repository

import "github.com/RuanScherer/journey-track-api/domain/model"

type ProjectRepository interface {
	Register(project *model.Project) error
	Save(project *model.Project) error
	FindByMemberId(memberId string) ([]*model.Project, error)
	FindById(id string) (*model.Project, error)
	FindMembersCountAndEventsCountById(id string) (*ProjectMembersCountAndEventsCount, error)
	DeleteById(id string) error
	HasMember(projectID, memberID string) (bool, error)
}

type ProjectMembersCountAndEventsCount struct {
	MembersCount int64
	EventsCount  int64
}
