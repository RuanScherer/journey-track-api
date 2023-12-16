package model

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

type User struct {
	ID                 string           `json:"id" valid:"uuid~[user] Invalid ID"`
	Email              string           `json:"email" valid:"required~[user] Email is required,email~[user] Invalid email"`
	Name               string           `json:"name" valid:"required~[user] Name is required,minstringlength(2)~[user] Name too short"`
	Password           string           `json:"password" valid:"required~[user] Password is required,minstringlength(8)~[user] Password too short"`
	VerificationToken  string           `valid:"-"`
	IsVerified         bool             `json:"is_verified" valid:"-"`
	PasswordResetToken string           `valid:"-"`
	Projects           []*Project       `json:"projects" valid:"-"`
	ProjectInvites     []*ProjectInvite `json:"project_invites" valid:"-"`
}

func NewUser(email string, name string, password string) (*User, error) {
	user := &User{
		ID:       uuid.New().String(),
		Email:    email,
		Name:     name,
		Password: password,
	}
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return nil, err
	}

	verificationToken := uuid.New().String()
	user.VerificationToken = verificationToken
	user.IsVerified = false

	return user, nil
}

func (user *User) RegenerateVerificationToken() error {
	if user.IsVerified {
		return errors.New("user already verified")
	}

	verificationToken := uuid.New().String()
	user.VerificationToken = verificationToken
	return nil
}

func (user *User) Verify(verificationToken string) error {
	if user.VerificationToken != verificationToken {
		return errors.New("invalid verification token")
	}

	user.IsVerified = true
	user.VerificationToken = ""
	return nil
}

func (user *User) ChangeName(newName string) error {
	user.Name = newName
	_, err := govalidator.ValidateStruct(user)
	return err
}

func (user *User) RequestPasswordReset() {
	passwordResetToken := uuid.New().String()
	user.PasswordResetToken = passwordResetToken
}

func (user *User) ResetPassword(newPassword string, passwordResetToken string) error {
	if user.PasswordResetToken != passwordResetToken {
		return errors.New("invalid password reset token")
	}

	user.Password = newPassword
	user.PasswordResetToken = ""

	_, err := govalidator.ValidateStruct(user)
	return err
}

func (user *User) HasPendingInviteForProject(project *Project) bool {
	for _, invite := range user.ProjectInvites {
		if invite.Project.ID == project.ID && invite.Status == ProjectInviteStatusPending {
			return true
		}
	}
	return false
}
