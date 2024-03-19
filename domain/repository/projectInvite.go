package repository

import "github.com/RuanScherer/journey-track-api/domain/model"

type ProjectInviteRepository interface {
	Create(projectInvite *model.ProjectInvite) error
	BatchCreate(projectInvites []*model.ProjectInvite) error
	Save(projectInvite *model.ProjectInvite) error
	DeleteById(projectInviteId string) error
	FindById(projectInviteId string) (*model.ProjectInvite, error)
	ListByProjectAndStatus(projectId string, status string) ([]*model.ProjectInvite, error)
	FindByProjectAndToken(projectId string, token string) (*model.ProjectInvite, error)
	FindPendingByUserAndProject(userId string, projectId string) (*model.ProjectInvite, error)
}
