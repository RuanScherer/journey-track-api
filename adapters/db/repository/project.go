package repository

import (
	"github.com/RuanScherer/journey-track-api/domain/model"
	"gorm.io/gorm"
)

type ProjectDBRepository struct {
	DB *gorm.DB
}

func NewProjectDBRepository(db *gorm.DB) *ProjectDBRepository {
	return &ProjectDBRepository{DB: db}
}

func (repository *ProjectDBRepository) Register(project *model.Project) error {
	return repository.DB.Create(project).Error
}

func (repository *ProjectDBRepository) Save(project *model.Project) error {
	return repository.DB.Save(project).Error
}

func (repository *ProjectDBRepository) FindByMemberId(memberId string) ([]*model.Project, error) {
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

func (repository *ProjectDBRepository) FindById(id string) (*model.Project, error) {
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

func (repository *ProjectDBRepository) DeleteById(id string) error {
	err := repository.DB.Where("id = ?", id).Delete(&model.Project{}).Error
	return err
}
