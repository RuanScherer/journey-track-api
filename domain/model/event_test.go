package model

import "testing"

func TestNewEvent(t *testing.T) {
	t.Run("should get error when project is invalid", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(*projectOwner.VerificationToken)
		project, err := NewProject("", projectOwner)
		if err == nil {
			t.Errorf("project should be invalid")
		}

		_, err = NewEvent("Event", project)
		if err == nil {
			t.Errorf("should get error when project is invalid")
		}
	})

	t.Run("should get error when provided name is invalid", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("Test", projectOwner)

		_, err := NewEvent("", project)
		if err == nil {
			t.Errorf("should get error when provided name is empty")
		}

		_, err = NewEvent("A", project)
		if err == nil {
			t.Errorf("should get error when provided name is too short")
		}
	})

	t.Run("should create event", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("Test", projectOwner)

		event, err := NewEvent("Test", project)
		if err != nil {
			t.Errorf("should create event")
		}

		if event.Timestamp == nil {
			t.Errorf("should set timestamp after creating event")
		}
	})
}
