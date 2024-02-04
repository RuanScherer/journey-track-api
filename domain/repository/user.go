package repository

import "github.com/RuanScherer/journey-track-api/domain/model"

type UserRepository interface {
	Register(user *model.User) error
	Save(user *model.User) error
	FindById(id string) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	Search(options UserSearchOptions) ([]*model.User, error)
}

type UserSearchOptions struct {
	Email    string
	Page     int
	PageSize int
}
