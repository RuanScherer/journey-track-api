package repository

import (
	"github.com/RuanScherer/journey-track-api/adapters/db/utils"
	"github.com/RuanScherer/journey-track-api/domain/model"
	domainrepository "github.com/RuanScherer/journey-track-api/domain/repository"
	"gorm.io/gorm"
)

type UserDBRepository struct {
	DB *gorm.DB
}

func NewUserDBRepository(db *gorm.DB) *UserDBRepository {
	return &UserDBRepository{DB: db}
}

func (repository *UserDBRepository) Register(user *model.User) error {
	return repository.DB.Create(user).Error
}

func (repository *UserDBRepository) Save(user *model.User) error {
	return repository.DB.Save(user).Error
}

func (repository *UserDBRepository) FindById(id string) (*model.User, error) {
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

func (repository *UserDBRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := repository.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repository *UserDBRepository) Search(options domainrepository.UserSearchOptions) ([]*model.User, error) {
	users := []*model.User{}
	err := repository.DB.
		Where("email like ?", "%"+options.Email+"%").
		Scopes(
			utils.Paginate(utils.PaginationOptions{
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
