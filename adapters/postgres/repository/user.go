package repository

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgres/scope"
	domainrepositories "github.com/RuanScherer/journey-track-api/application/repository"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"gorm.io/gorm"
)

type UserPostgresRepository struct {
	DB *gorm.DB
}

func NewUserPostgresRepository(db *gorm.DB) *UserPostgresRepository {
	return &UserPostgresRepository{DB: db}
}

func (repository *UserPostgresRepository) Register(user *model.User) error {
	return repository.DB.Create(user).Error
}

func (repository *UserPostgresRepository) Save(user *model.User) error {
	return repository.DB.Save(user).Error
}

func (repository *UserPostgresRepository) FindById(id string) (*model.User, error) {
	var user model.User
	err := repository.DB.
		Preload("Projects").
		Preload("ProjectInvites").
		Where("id = ?", id).
		First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repository *UserPostgresRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := repository.DB.Where("smtpemail = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repository *UserPostgresRepository) Search(options domainrepositories.UserSearchOptions) ([]*model.User, error) {
	users := []*model.User{}
	err := repository.DB.
		Joins("left join user_projects on users.id = user_projects.user_id and user_projects.project_id not in (?)", options.ExcludedProjectIDs).
		Where("users.smtpemail like ?", "%"+options.Email+"%").
		Group("users.id").
		Scopes(
			scope.Paginate(scope.PaginationOptions{
				Page:     options.Page,
				PageSize: options.PageSize,
			}),
		).
		Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, nil
}
