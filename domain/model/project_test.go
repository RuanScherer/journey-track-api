package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProject(t *testing.T) {
	t.Run("should get error when owner is invalid", func(t *testing.T) {
		_, err := NewProject("project name", &User{})
		require.NotNil(t, err)
	})

	t.Run("should get error when owner is not verified", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")

		_, err := NewProject("project name", owner)
		require.NotNil(t, err)
		require.Equal(t, "owner must be verified", err.Error())
	})

	t.Run("should get error when project name is invalid", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(*owner.VerificationToken)

		_, err := NewProject("", owner)
		require.NotNil(t, err)
		require.Equal(t, "[project] Name is required", err.Error())

		_, err = NewProject("a", owner)
		require.NotNil(t, err)
		require.Equal(t, "[project] Name too short", err.Error())
	})

	t.Run("should get project when owner and project name is valid", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(*owner.VerificationToken)

		project, err := NewProject("my project", owner)
		require.Nil(t, err)
		require.NotNil(t, project)
		require.NotNil(t, project.Token)
		require.NotEmpty(t, *project.Token)
	})
}

func TestProjectChangeName(t *testing.T) {
	t.Run("should get error when provided name is invalid", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(*owner.VerificationToken)

		project, _ := NewProject("my project", owner)
		err := project.ChangeName("")
		require.NotNil(t, err)
		require.Equal(t, "[project] Name is required", err.Error())

		err = project.ChangeName("a")
		require.NotNil(t, err)
		require.Equal(t, "[project] Name too short", err.Error())
	})

	t.Run("should change name when provided name is valid", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(*owner.VerificationToken)

		project, _ := NewProject("my project", owner)
		project.ChangeName("my other project")
		require.Equal(t, "my other project", project.Name)
	})
}

func TestAddMember(t *testing.T) {
	t.Run("should get error when provided user is invalid", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(*owner.VerificationToken)

		project, _ := NewProject("my project", owner)

		newMember, _ := NewUser("", "Jon Doe", "pass1234")
		err := project.AddMember(newMember)
		require.NotNil(t, err)
	})

	t.Run("should get error when provided user is already a member of the project", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(*owner.VerificationToken)

		project, _ := NewProject("my project", owner)

		newMember, _ := NewUser("jondoe@test.com", "Jon Doe", "pass1234")
		project.AddMember(newMember)

		err := project.AddMember(newMember)
		require.NotNil(t, err)
		require.Equal(t, "user is already a member of the project", err.Error())
	})

	t.Run("should add member", func(t *testing.T) {
		owner, _ := NewUser("owner@domain.com", "Owner", "pass1234")
		owner.Verify(*owner.VerificationToken)

		project, _ := NewProject("my project", owner)

		newMember, _ := NewUser("jondoe@test.com", "Jon Doe", "pass1234")
		err := project.AddMember(newMember)
		require.Nil(t, err)
		require.Len(t, project.Members, 2)
	})
}
