package model

import (
	"testing"
)

func TestNewProject(t *testing.T) {
	t.Run("should get error when owner is invalid", func(t *testing.T) {
		_, err := NewProject("project name", &User{})
		if err == nil {
			t.Error("should get error when owner is invalid")
		}
	})

	t.Run("should get error when owner is not verified", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")

		_, err := NewProject("project name", owner)
		if err == nil {
			t.Error("should get error when owner is not verified")
		}
	})

	t.Run("should get error when project name is invalid", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(owner.VerificationToken)

		_, err := NewProject("", owner)
		if err == nil {
			t.Error("should get error when project name is empty")
		}

		_, err = NewProject("a", owner)
		if err == nil {
			t.Error("should get error when project name is too short")
		}
	})

	t.Run("should get project when owner and project name is valid", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(owner.VerificationToken)

		_, err := NewProject("my project", owner)
		if err != nil {
			t.Error("should get project when owner and project name is valid")
		}
	})

	t.Run("should have token", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(owner.VerificationToken)

		project, _ := NewProject("my project", owner)
		if project.Token == "" {
			t.Error("should have token")
		}
	})

	t.Run("owner should have project", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(owner.VerificationToken)

		project, _ := NewProject("my project", owner)
		if owner.Projects[0].ID != project.ID {
			t.Error("owner should have project")
		}
	})
}

func TestProjectChangeName(t *testing.T) {
	t.Run("should get error when provided name is invalid", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(owner.VerificationToken)

		project, _ := NewProject("my project", owner)
		err := project.ChangeName("")

		if err == nil {
			t.Error("should get error when provided name is empty")
		}

		err = project.ChangeName("a")
		if err == nil {
			t.Error("should get error when provided name is too short")
		}
	})

	t.Run("should change name when provided name is valid", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(owner.VerificationToken)

		project, _ := NewProject("my project", owner)
		project.ChangeName("my other project")

		if project.Name != "my other project" {
			t.Error("should change name when provided name is valid")
		}
	})
}
