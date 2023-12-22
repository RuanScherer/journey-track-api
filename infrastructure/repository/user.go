package repository

import (
	"github.com/RuanScherer/journey-track-api/domain/model"
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
	return &user, err
}

func (repository *UserDBRepository) SearchByEmail(email string) ([]*model.User, error) {
	users := []*model.User{}
	err := repository.DB.Where("email like ?", "%"+email+"%").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
