package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	t.Run("should return error when email is invalid", func(t *testing.T) {
		_, err := NewUser("", "name", "password")
		require.NotNil(t, err)
		require.Equal(t, "[user] Email is required", err.Error())

		_, err = NewUser("invalid email", "name", "password")
		require.NotNil(t, err)
		require.Equal(t, "[user] Invalid email", err.Error())

		_, err = NewUser("invalid@invalid", "name", "password")
		require.NotNil(t, err)
		require.Equal(t, "[user] Invalid email", err.Error())
	})

	t.Run("should return error when name is invalid", func(t *testing.T) {
		_, err := NewUser("example@domain.com", "", "password")
		require.NotNil(t, err)
		require.Equal(t, "[user] Name is required", err.Error())

		_, err = NewUser("example@domain.com", "n", "password")
		require.NotNil(t, err)
		require.Equal(t, "[user] Name too short", err.Error())
	})

	t.Run("should return error when password is invalid", func(t *testing.T) {
		_, err := NewUser("example@domain.com", "Ruan", "")
		require.NotNil(t, err)
		require.Equal(t, "[user] Password is required", err.Error())

		_, err = NewUser("example@domain.com", "Ruan", "1234567")
		require.NotNil(t, err)
		require.Equal(t, "[user] Password too short", err.Error())
	})

	t.Run("should return user when data is valid", func(t *testing.T) {
		user, err := NewUser("example@domain.com", "Ruan", "12345678")
		require.Nil(t, err)
		require.NotNil(t, user)
		require.NotEmpty(t, user.ID)
		require.NotNil(t, user.VerificationToken)
		require.NotEmpty(t, *user.VerificationToken)
		require.False(t, user.IsVerified)
	})
}

func TestRegenerateVerificationToken(t *testing.T) {
	t.Run("should not regenerate verification token when user is already verified", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		user.Verify(*user.VerificationToken)

		err := user.RegenerateVerificationToken()
		require.NotNil(t, err)
		require.Equal(t, "user already verified", err.Error())
	})

	t.Run("should regenerate verification token when user is not verified", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		originalVerificationToken := user.VerificationToken

		err := user.RegenerateVerificationToken()
		require.Nil(t, err)
		require.NotEqual(t, originalVerificationToken, user.VerificationToken)
		require.NotNil(t, user.VerificationToken)
		require.NotEmpty(t, *user.VerificationToken)
	})
}

func TestVerify(t *testing.T) {
	t.Run("should return error when verification token is invalid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		err := user.Verify("oisnoanfiosanf")
		require.NotNil(t, err)
		require.Equal(t, "invalid verification token", err.Error())
	})

	t.Run("should verify user when verification token is valid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		err := user.Verify(*user.VerificationToken)
		require.Nil(t, err)
		require.True(t, user.IsVerified)
		require.Nil(t, user.VerificationToken)
	})
}

func TestUserChangeName(t *testing.T) {
	t.Run("should return error when provided name is invalid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")

		err := user.ChangeName("")
		require.NotNil(t, err)
		require.Equal(t, "[user] Name is required", err.Error())

		err = user.ChangeName("n")
		require.NotNil(t, err)
		require.Equal(t, "[user] Name too short", err.Error())
	})

	t.Run("should change name when provided name is valid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")

		err := user.ChangeName("Ruan Scherer")
		require.Nil(t, err)
		require.Equal(t, "Ruan Scherer", user.Name)
	})
}

func TestRequestPasswordReset(t *testing.T) {
	user, _ := NewUser("example@domain.com", "Ruan", "12345678")
	require.Nil(t, user.PasswordResetToken)

	user.RequestPasswordReset()
	require.NotNil(t, user.PasswordResetToken)
	require.NotEmpty(t, *user.PasswordResetToken)
}

func TestResetPassword(t *testing.T) {
	t.Run("should return error when there's no request for password reset", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")

		err := user.ResetPassword("new password", "")
		require.NotNil(t, err)
		require.Equal(t, "user has no request for password reset", err.Error())
	})

	t.Run("should return error when password reset token is invalid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		user.RequestPasswordReset()

		err := user.ResetPassword("new password", "invalid token")
		require.NotNil(t, err)
		require.Equal(t, "invalid password reset token", err.Error())
	})

	t.Run("should get error when provided password is invalid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		user.RequestPasswordReset()

		err := user.ResetPassword("", *user.PasswordResetToken)
		require.NotNil(t, err)
		require.Equal(t, "[user] Password is required", err.Error())

		err = user.ResetPassword("1234567", *user.PasswordResetToken)
		require.NotNil(t, err)
		require.Equal(t, "[user] Password too short", err.Error())
	})

	t.Run("should reset password when password reset token is valid", func(t *testing.T) {
		user, _ := NewUser("example@domain.com", "Ruan", "12345678")
		user.RequestPasswordReset()

		err := user.ResetPassword("new password", *user.PasswordResetToken)
		require.Nil(t, err)
	})
}
