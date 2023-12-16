package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

const (
	ProjectInviteStatusPending  = "pending"
	ProjectInviteStatusAccepted = "accepted"
	ProjectInviteStatusDeclined = "declined"
	ProjectInviteStatusRevoked  = "revoked"
)

type ProjectInvite struct {
	ID              string     `json:"id" valid:"uuid~[project invite] Invalid ID"`
	Project         *Project   `json:"project" valid:"-"`
	User            *User      `json:"user" valid:"-"`
	Status          string     `json:"status" valid:"in(pending|accepted|declined|revoked)~[project invite] Invalid status"`
	Token           string     `valid:"uuid~[project invite] Invalid token"`
	AnswerTimestamp *time.Time `json:"answer_timestamp" valid:"-"`
	RevokeTimestamp *time.Time `json:"revoke_timestamp" valid:"-"`
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

	if user.HasPendingInviteForProject(project) {
		return nil, errors.New("user already has a pending invite for the project")
	}

	projectInvite := &ProjectInvite{
		ID:      uuid.New().String(),
		Project: project,
		User:    user,
		Status:  ProjectInviteStatusPending,
		Token:   uuid.New().String(),
	}

	_, err = govalidator.ValidateStruct(projectInvite)
	if err != nil {
		return nil, err
	}

	user.ProjectInvites = append(user.ProjectInvites, projectInvite)
	project.Invites = append(project.Invites, projectInvite)
	return projectInvite, nil
}

func (projectInvite *ProjectInvite) Accept(token string) error {
	err := projectInvite.answer(ProjectInviteStatusAccepted, token)
	if err != nil {
		return err
	}

	projectInvite.Project.Members = append(projectInvite.Project.Members, projectInvite.User)
	return nil
}

func (projectInvite *ProjectInvite) Decline(token string) error {
	err := projectInvite.answer(ProjectInviteStatusDeclined, token)
	return err
}

func (projectInvite *ProjectInvite) answer(answer string, token string) error {
	if projectInvite.Token != token {
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
	answerTimestamp := time.Now()
	projectInvite.AnswerTimestamp = &answerTimestamp

	_, err := govalidator.ValidateStruct(projectInvite)
	return err
}

func (projectInvite *ProjectInvite) Revoke(actor *User) error {
	if !projectInvite.Project.HasMember(actor) {
		return errors.New("only project members can revoke invites")
	}

	if projectInvite.Status != ProjectInviteStatusPending {
		return errors.New("invite already answered or revoked")
	}

	projectInvite.Status = ProjectInviteStatusRevoked
	revokeTimestamp := time.Now()
	projectInvite.RevokeTimestamp = &revokeTimestamp

	_, err := govalidator.ValidateStruct(projectInvite)
	return err
}
