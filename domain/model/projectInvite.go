package model

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	ProjectInviteStatusPending  = "pending"
	ProjectInviteStatusAccepted = "accepted"
	ProjectInviteStatusDeclined = "declined"
)

type ProjectInviteRepository interface {
	Create(projectInvite *ProjectInvite) error
	Save(projectInvite *ProjectInvite) error
	DeleteById(projectInviteId string) error
	FindById(projectInviteId string) (*ProjectInvite, error)
	FindByProjectAndToken(projectId string, token string) (*ProjectInvite, error)
	FindPendingByUserAndProject(userId string, projectId string) (*ProjectInvite, error)
}

type ProjectInvite struct {
	gorm.Model
	ID        string   `json:"id" gorm:"primaryKey" valid:"uuid~[project invite] Invalid ID"`
	ProjectID string   `gorm:"column:project_id;type:varchar(255);not null" valid:"-"`
	Project   *Project `json:"project" gorm:"" valid:"-"`
	UserID    string   `gorm:"column:user_id;type:varchar(255);not null" valid:"-"`
	User      *User    `json:"user" valid:"-"`
	Status    string   `json:"status" gorm:"type:varchar(100);not null" valid:"in(pending|accepted|declined|revoked)~[project invite] Invalid status"`
	Token     *string  `gorm:"type:varchar(255);unique;not null" valid:"uuid~[project invite] Invalid token"`
}

func NewProjectInvite(project *Project, user *User) (*ProjectInvite, error) {
	_, err := govalidator.ValidateStruct(project)
	if err != nil {
		return nil, err
	}

	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		return nil, err
	}

	if project.HasMember(user) {
		return nil, errors.New("user is already a member of the project")
	}

	token := uuid.New().String()
	projectInvite := &ProjectInvite{
		ID:      uuid.New().String(),
		Project: project,
		UserID:  user.ID,
		Status:  ProjectInviteStatusPending,
		Token:   &token,
	}

	_, err = govalidator.ValidateStruct(projectInvite)
	if err != nil {
		return nil, err
	}

	return projectInvite, nil
}

func (projectInvite *ProjectInvite) Accept(token string) error {
	err := projectInvite.answer(ProjectInviteStatusAccepted, token)
	return err
}

func (projectInvite *ProjectInvite) Decline(token string) error {
	err := projectInvite.answer(ProjectInviteStatusDeclined, token)
	return err
}

func (projectInvite *ProjectInvite) answer(answer string, token string) error {
	if *projectInvite.Token != token {
		return errors.New("invalid token provided to answer invite")
	}

	if projectInvite.Status != ProjectInviteStatusPending {
		return errors.New("invite already answered or revoked")
	}

	isValidAnswer := govalidator.IsIn(answer, ProjectInviteStatusAccepted, ProjectInviteStatusDeclined)
	if !isValidAnswer {
		return errors.New("invalid answer provided to invite")
	}

	projectInvite.Status = answer
	_, err := govalidator.ValidateStruct(projectInvite)
	return err
}
