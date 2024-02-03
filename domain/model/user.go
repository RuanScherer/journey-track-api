package model

import (
	"errors"

	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                 string           `json:"id" gorm:"primaryKey" valid:"uuid~[user] Invalid ID"`
	Email              *string          `json:"email" gorm:"type:varchar(255);unique;not null" valid:"required~[user] Email is required,email~[user] Invalid email"`
	Name               string           `json:"name" gorm:"type:varchar(255);not null" valid:"required~[user] Name is required,minstringlength(2)~[user] Name too short"`
	Password           string           `json:"password" gorm:"type:varchar(255);not null" valid:"required~[user] Password is required,minstringlength(8)~[user] Password too short"`
	VerificationToken  *string          `gorm:"column:verification_token;type:varchar(255);unique;default:null" valid:"-"`
	IsVerified         bool             `json:"is_verified" gorm:"column:is_verified;type:boolean;not null" valid:"-"`
	PasswordResetToken *string          `gorm:"column:password_reset_token;type:varchar(255);unique;default:null" valid:"-"`
	Projects           []*Project       `gorm:"many2many:user_projects" json:"projects" valid:"-"`
	ProjectInvites     []*ProjectInvite `gorm:"foreignKey:UserID" json:"project_invites" valid:"-"`
}

func NewUser(email string, name string, password string) (*User, error) {
	user := &User{
		ID:       uuid.New().String(),
		Email:    &email,
		Name:     name,
		Password: password,
	}
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("[user] Unable to encrypt password")
	}
	user.Password = string(passwordHash)

	verificationToken := uuid.New().String()
	user.VerificationToken = &verificationToken
	user.IsVerified = false

	return user, nil
}

func (user *User) RegenerateVerificationToken() error {
	if user.IsVerified {
		return errors.New("user already verified")
	}

	verificationToken := uuid.New().String()
	user.VerificationToken = &verificationToken
	return nil
}

func (user *User) Verify(verificationToken string) error {
	if user.IsVerified {
		return errors.New("user already verified")
	}

	if user.VerificationToken == nil || *user.VerificationToken != verificationToken {
		return errors.New("invalid verification token")
	}

	user.IsVerified = true
	user.VerificationToken = nil
	return nil
}

func (user *User) ChangeName(newName string) error {
	user.Name = newName
	_, err := govalidator.ValidateStruct(user)
	return err
}

func (user *User) RequestPasswordReset() {
	passwordResetToken := uuid.New().String()
	user.PasswordResetToken = &passwordResetToken
}

func (user *User) ResetPassword(newPassword string, passwordResetToken string) error {
	if user.PasswordResetToken == nil {
		return errors.New("user has no request for password reset")
	}

	if *user.PasswordResetToken != passwordResetToken {
		return errors.New("invalid password reset token")
	}

	user.Password = newPassword
	_, err := govalidator.ValidateStruct(user)
	if err != nil {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("[user] Unable to encrypt password")
	}
	user.Password = string(passwordHash)

	user.PasswordResetToken = nil
	return nil
}
