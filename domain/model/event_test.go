package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewEvent(t *testing.T) {
	t.Run("should get error when project is invalid", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, err := NewProject("", projectOwner)
		require.NotNil(t, err)

		_, err = NewEvent("Event", project)
		require.NotNil(t, err)
	})

	t.Run("should get error when provided name is invalid", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("Test", projectOwner)

		_, err := NewEvent("", project)
		require.NotNil(t, err)
		require.Equal(t, "[event] Name is required", err.Error())

		_, err = NewEvent("A", project)
		require.NotNil(t, err)
		require.Equal(t, "[event] Name should be longer than 2 characters", err.Error())
	})

	t.Run("should create event", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("Test", projectOwner)

		event, err := NewEvent("Test", project)
		require.Nil(t, err)
		require.NotNil(t, event)
		require.NotNil(t, event.Timestamp)
	})
}
