package model

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	t.Run("should return error when email is invalid", func(t *testing.T) {
		_, err := NewUser("", "name", "password")
		if err == nil || err.Error() != "[user] Email is required" {
			t.Error("should return error when email is blank")
		}

		_, err = NewUser("invalid email", "name", "password")
		if err == nil || err.Error() != "[user] Invalid email" {
			t.Error("should return error when email is invalid")
		}

		_, err = NewUser("invalid@invalid", "name", "password")
		if err == nil || err.Error() != "[user] Invalid email" {
			t.Error("should return error when email is invalid")
		}
	})

	t.Run("should return error when name is invalid", func(t *testing.T) {
		_, err := NewUser("example@domain.com", "", "password")
		if err == nil || err.Error() != "[user] Name is required" {
			t.Error("should return error when name is blank")
		}

		_, err = NewUser("example@domain.com", "n", "password")
		if err == nil || err.Error() != "[user] Name too short" {
			t.Error("should return error when name is too short")
		}
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		_, err := NewUser("example@domain.com", "Ruan", "")
		if err == nil || err.Error() != "[user] Password is required" {
			t.Error("should return error when password is blank")
		}

		_, err = NewUser("example@domain.com", "Ruan", "1234567")
		if err == nil || err.Error() != "[user] Password too short" {
			t.Error("should return error when password is too short")
		}
	})

	t.Run("should return user when data is valid", func(t *testing.T) {
		user, err := NewUser("example@domain.com", "Ruan", "12345678")
		if err != nil {
			t.Error("should return user when data is valid")
			return
		}

		if user == nil {
			t.Error("should return user when data is valid")
			return
		}

		if user.ID == "" {
			t.Error("should return user with ID")
		}
	})

	t.Run("should return user with verification token", func(t *testing.T) {
		user, err := NewUser("example@domain.com", "Ruan", "12345678")
		if err != nil {
			t.Error("should return user when data is valid")
			return
		}

		if user == nil {
			t.Error("should return user when data is valid")
			return
		}

		if user.VerificationToken == "" {
			t.Error("should return user with verification token")
		}
	})

	t.Run("should return user not verified", func(t *testing.T) {
		user, err := NewUser("example@domain.com", "Ruan", "12345678")
		if err != nil {
			t.Error("should return user when data is valid")
			return
		}

		if user == nil {
			t.Error("should return user when data is valid")
			return
		}

		if user.IsVerified {
			t.Error("should return user not verified")
		}
	})
}

func TestRegenerateVerificationToken(t *testing.T) {
	t.Run("should not regenerate verification token when user is already verified", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		user.Verify(user.VerificationToken)

		err := user.RegenerateVerificationToken()
		if err == nil || err.Error() != "user already verified" {
			t.Error("should not regenerate verification token when user is already verified")
		}
	})

	t.Run("should regenerate verification token when user is not verified", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		originalVerificationToken := user.VerificationToken

		err := user.RegenerateVerificationToken()
		if err != nil {
			t.Error("should not return error when user is not verified")
		}

		if user.VerificationToken == originalVerificationToken {
			t.Error("should regenerate verification token when user is not verified")
		}
	})
}

func TestVerify(t *testing.T) {
	t.Run("should return error when verification token is invalid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		err := user.Verify("oisnoanfiosanf")

		if err == nil || err.Error() != "invalid verification token" {
			t.Error("should return error when verification token is invalid")
		}
	})

	t.Run("should verify user when verification token is valid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		err := user.Verify(user.VerificationToken)

		if err != nil {
			t.Error("should not get error when verification token is valid")
		}

		if !user.IsVerified {
			t.Error("should verify user when verification token is valid")
		}

		if user.VerificationToken != "" {
			t.Error("should delete verification token when already used")
		}
	})
}

func TestUserChangeName(t *testing.T) {
	t.Run("should return error when provided name is invalid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")

		err := user.ChangeName("")
		if err == nil || err.Error() != "[user] Name is required" {
			t.Error("should return error when name is blank")
		}

		err = user.ChangeName("n")
		if err == nil || err.Error() != "[user] Name too short" {
			t.Error("should return error when name is too short")
		}
	})
}

func TestRequestPasswordReset(t *testing.T) {
	user, _ := NewUser("example@domain.com", "Ruan", "12345678")
	if user.PasswordResetToken != "" {
		t.Error("should not have password reset token by default")
	}

	user.RequestPasswordReset()
	if user.PasswordResetToken == "" {
		t.Error("should have password reset token after request")
	}
}

func TestResetPassword(t *testing.T) {
	t.Run("should return error when password reset token is invalid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		user.RequestPasswordReset()

		err := user.ResetPassword("new password", "invalid token")
		if err == nil || err.Error() != "invalid password reset token" {
			t.Error("should return error when password reset token is invalid")
		}
	})

	t.Run("should get error when provided password is invalid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		user.RequestPasswordReset()

		err := user.ResetPassword("", user.PasswordResetToken)
		if err == nil || err.Error() != "[user] Password is required" {
			t.Error("should return error when password is blank")
		}

		err = user.ResetPassword("1234567", user.PasswordResetToken)
		if err == nil || err.Error() != "[user] Password too short" {
			t.Error("should return error when password is too short")
		}
	})

	t.Run("should reset password when password reset token is valid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		user.RequestPasswordReset()

		err := user.ResetPassword("new password", user.PasswordResetToken)
		if err != nil {
			t.Error("should reset password when password reset token is valid")
		}
	})
}

func TestHasPendingInviteForProject(t *testing.T) {
	t.Run("should return false when user has no pending invites for project", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")

		projectOwner, _ := NewUser("owner@domain.com", "Owner", "12345678")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		if user.HasPendingInviteForProject(project) {
			t.Error("should return false when user has no pending invites for project")
		}
	})

	t.Run("should return true when user has pending invites for project", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@domain.com", "Owner", "12345678")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		invitedUser, _ := NewUser("invited@domain.com", "Invited", "12345678")
		NewProjectInvite(project, invitedUser)

		if !invitedUser.HasPendingInviteForProject(project) {
			t.Error("should return faltruese when user has no pending invites for project")
		}
	})
}
