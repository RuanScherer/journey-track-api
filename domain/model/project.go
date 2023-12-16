package model

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type Project struct {
	ID      string           `json:"id" valid:"uuid~[project] Invalid ID"`
	Name    string           `json:"name" valid:"required~[project] Name is required,minstringlength(2)~[project] Name too short"`
	OwnerID string           `json:"owner_id" valid:"required~[project] Owner is required,uuid~[project] Invalid owner ID"`
	Members []*User          `json:"members" valid:"-"`
	Invites []*ProjectInvite `json:"invites" valid:"-"`
	Token   string           `json:"token" valid:"required~[project] Token is required,uuid~[project] Invalid token"`
}

func NewProject(name string, owner *User) (*Project, error) {
	_, err := govalidator.ValidateStruct(owner)
	if err != nil {
		return nil, err
	}

	if !owner.IsVerified {
		return nil, errors.New("owner must be verified")
	}

	project := &Project{
		ID:      uuid.New().String(),
		Name:    name,
		OwnerID: owner.ID,
		Members: []*User{owner},
		Token:   uuid.New().String(),
	}

	_, err = govalidator.ValidateStruct(project)
	if err != nil {
		return nil, err
	}

	owner.Projects = append(owner.Projects, project)
	return project, nil
}

func (project *Project) ChangeName(newName string) error {
	project.Name = newName
	_, err := govalidator.ValidateStruct(project)
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
