package repositories

import (
	"github.com/RuanScherer/journey-track-api/domain/model"
	domainrepositories "github.com/RuanScherer/journey-track-api/domain/repository"
	"gorm.io/gorm"
)

type ProjectPostgresRepository struct {
	DB *gorm.DB
}

func NewProjectPostgresRepository(db *gorm.DB) *ProjectPostgresRepository {
	return &ProjectPostgresRepository{DB: db}
}

func (repository *ProjectPostgresRepository) Register(project *model.Project) error {
	return repository.DB.Create(project).Error
}

func (repository *ProjectPostgresRepository) Save(project *model.Project) error {
	return repository.DB.Save(project).Error
}

func (repository *ProjectPostgresRepository) FindByMemberId(memberId string) ([]*model.Project, error) {
	projects := []*model.Project{}
	err := repository.DB.
		Joins("inner join user_projects on user_projects.project_id = projects.id").
		Where("user_projects.user_id = ?", memberId).
		Find(&projects).Error

	if err != nil {
		return nil, err
	}
	return projects, nil
}

func (repository *ProjectPostgresRepository) FindById(id string) (*model.Project, error) {
	project := &model.Project{}
	err := repository.DB.
		Preload("Members").
		Preload("Invites").
		Where("id = ?", id).
		First(project).Error

	if err != nil {
		return nil, err
	}
	return project, nil
}

func (repository *ProjectPostgresRepository) FindMembersCountAndEventsCountById(
	id string,
) (*domainrepositories.ProjectInvitesCountAndEventsCount, error) {
	result := &domainrepositories.ProjectInvitesCountAndEventsCount{}
	err := repository.DB.
		Table("projects").
		Joins("left join project_invites on projects.id = project_invites.project_id").
		Joins("left join events on projects.id = events.project_id").
		Where("projects.id = ? and project_invites.deleted_at is null", id).
		Select("count(project_invites.id) as invites_count, count(events.id) as events_count").
		Scan(result).Error

	if err != nil {
		return nil, err
	}
	return result, nil
}

func (repository *ProjectPostgresRepository) DeleteById(id string) error {
	err := repository.DB.Where("id = ?", id).Delete(&model.Project{}).Error
	return err
}

func (repository *ProjectPostgresRepository) HasMember(projectID, memberID string) (bool, error) {
	var count int64
	err := repository.DB.
		Table("user_projects").
		Where("project_id = ?", projectID).
		Where("user_id = ?", memberID).
		Limit(1).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}
