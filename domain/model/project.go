package model

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository interface {
	Register(project *Project) error
	Save(project *Project) error
	FindByMemberId(memberId string) ([]*Project, error)
	FindById(id string) (*Project, error)
	DeleteById(id string) error
}

type Project struct {
	gorm.Model
	ID      string           `json:"id" gorm:"primaryKey" valid:"uuid~[project] Invalid ID"`
	Name    string           `json:"name" gorm:"type:varchar(255);not null" valid:"required~[project] Name is required,minstringlength(2)~[project] Name too short"`
	OwnerID string           `json:"owner_id" gorm:"column:owner_id;type:varchar(255);not null" valid:"required~[project] Owner is required,uuid~[project] Invalid owner ID"`
	Members []*User          `json:"members" gorm:"many2many:user_projects" valid:"-"`
	Invites []*ProjectInvite `json:"invites" gorm:"foreignKey:ProjectID" valid:"-"`
	Events  []*Event         `json:"events" gorm:"foreignKey:ProjectID" valid:"-"`
	Token   *string          `json:"token" gorm:"type:varchar(255);unique" valid:"required~[project] Token is required,uuid~[project] Invalid token"`
}

func NewProject(name string, owner *User) (*Project, error) {
	_, err := govalidator.ValidateStruct(owner)
	if err != nil {
		return nil, err
	}

	if !owner.IsVerified {
		return nil, errors.New("owner must be verified")
	}

	token := uuid.New().String()
	project := &Project{
		ID:      uuid.New().String(),
		Name:    name,
		OwnerID: owner.ID,
		Members: []*User{owner},
		Token:   &token,
	}

	_, err = govalidator.ValidateStruct(project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

func (project *Project) ChangeName(newName string) error {
	project.Name = newName
	_, err := govalidator.ValidateStruct(project)
	return err
}

func (project *Project) AddMember(user *User) error {
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return err
	}

	if project.HasMember(user) {
		return errors.New("user is already a member of the project")
	}

	project.Members = append(project.Members, user)
	_, err = govalidator.ValidateStruct(project)
	return err
}

func (project *Project) HasMember(user *User) bool {
	for _, member := range project.Members {
		if member.ID == user.ID {
			return true
		}
	}
	return false
}
