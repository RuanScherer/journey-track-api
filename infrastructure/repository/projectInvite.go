package repository

import (
	"github.com/RuanScherer/journey-track-api/domain/model"
	"gorm.io/gorm"
)

type ProjectInviteDBRepository struct {
	DB *gorm.DB
}

func NewProjectInviteDBRepository(db *gorm.DB) *ProjectInviteDBRepository {
	return &ProjectInviteDBRepository{DB: db}
}

func (repository *ProjectInviteDBRepository) Create(projectInvite *model.ProjectInvite) error {
	return repository.DB.Create(projectInvite).Error
}

func (repository *ProjectInviteDBRepository) Save(projectInvite *model.ProjectInvite) error {
	return repository.DB.Save(projectInvite).Error
}

func (repository *ProjectInviteDBRepository) DeleteById(projectInviteId string) error {
	err := repository.DB.Where("id = ?", projectInviteId).Delete(&model.ProjectInvite{}).Error
	return err
}

func (repository *ProjectInviteDBRepository) FindById(projectInviteId string) (*model.ProjectInvite, error) {
	projectInvite := &model.ProjectInvite{}
	err := repository.DB.
		Preload("User").
		Preload("Project.Members").
		Where("id = ?", projectInviteId).
		First(projectInvite).Error

	if err != nil {
		return nil, err
	}
	return projectInvite, nil
}

func (repository *ProjectInviteDBRepository) FindByProjectAndToken(
	projectId string,
	token string,
) (*model.ProjectInvite, error) {
	projectInvite := &model.ProjectInvite{}
	err := repository.DB.
		Preload("User").
		Preload("Project").
		Where("project_id = ? and token = ?", projectId, token).
		First(projectInvite).Error

	if err != nil {
		return nil, err
	}
	return projectInvite, nil
}

func (repository *ProjectInviteDBRepository) FindPendingByUserAndProject(
	userId string,
	projectId string,
) (*model.ProjectInvite, error) {
	projectInvite := &model.ProjectInvite{}
	err := repository.DB.
		Where("user_id = ? and project_id = ? and status = ?", userId, projectId, model.ProjectInviteStatusPending).
		First(projectInvite).Error

	if err != nil {
		return nil, err
	}
	return projectInvite, nil
}
