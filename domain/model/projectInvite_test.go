package model

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewProjectInvite(t *testing.T) {
	t.Run("should get error when provided project is invalid", func(t *testing.T) {
		_, err := NewProjectInvite(&Project{}, &User{})
		require.NotNil(t, err)
	})

	t.Run("should get error when provided user is invalid", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		project, _ := NewProject("test", projectOwner)
		_, err := NewProjectInvite(project, &User{})
		require.NotNil(t, err)
	})

	t.Run("should get error when user is already a member", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, err := NewProjectInvite(project, userToInvite)
		require.Nil(t, err)

		_ = invite.Accept(*invite.Token)
		project.Members = append(project.Members, userToInvite)

		_, err = NewProjectInvite(project, userToInvite)
		require.NotNil(t, err)
		require.Equal(t, "user is already a member of the project", err.Error())
	})

	t.Run("should return project invite when provided project and user are valid", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, err := NewProjectInvite(project, userToInvite)
		require.Nil(t, err)
		require.Equal(t, ProjectInviteStatusPending, invite.Status)
		require.NotNil(t, invite.Token)
		require.NotEmpty(t, *invite.Token)
	})
}

func TestAccept(t *testing.T) {
	t.Run("should get error when provided token is invalid", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		err := invite.Accept("invalid-token")
		require.NotNil(t, err)
		require.Equal(t, "invalid token provided to answer invite", err.Error())
	})

	t.Run("should get error when invite is not pending", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)
		invite.Status = ProjectInviteStatusAccepted

		err := invite.Accept(*invite.Token)
		require.NotNil(t, err)
		require.Equal(t, "invite already answered or revoked", err.Error())
	})

	t.Run("should accept invite", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		err := invite.Accept(*invite.Token)
		require.Nil(t, err)
		require.Equal(t, ProjectInviteStatusAccepted, invite.Status)
	})
}

func TestDecline(t *testing.T) {
	t.Run("should get error when provided token is invalid", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		err := invite.Decline("invalid-token")
		require.NotNil(t, err)
		require.Equal(t, "invalid token provided to answer invite", err.Error())
	})

	t.Run("should get error when invited is not pending", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)
		invite.Status = ProjectInviteStatusAccepted

		err := invite.Decline(*invite.Token)
		require.NotNil(t, err)
		require.Equal(t, "invite already answered or revoked", err.Error())
	})

	t.Run("should decline invite", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		err := invite.Decline(*invite.Token)
		require.Nil(t, err)
		require.Equal(t, ProjectInviteStatusDeclined, invite.Status)
	})
}

func TestCanRevoke(t *testing.T) {
	t.Run("should not be able to revoke", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		invite.Status = ProjectInviteStatusAccepted
		canRevoke, reason := invite.CanRevoke()
		require.False(t, canRevoke)
		require.Equal(t, "invite already answered or revoked", reason)

		invite.Status = ProjectInviteStatusDeclined
		canRevoke, reason = invite.CanRevoke()
		require.False(t, canRevoke)
		require.Equal(t, "invite already answered or revoked", reason)
	})

	t.Run("should be able to revoke", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		canRevoke, reason := invite.CanRevoke()
		require.True(t, canRevoke)
		require.Empty(t, reason)
	})
}

func TestAnswer(t *testing.T) {
	t.Run("should get error when trying to use an invalid answer", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		_ = projectOwner.Verify(*projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		err := invite.answer("invalid-answer", *invite.Token)
		require.NotNil(t, err)
		require.Equal(t, "invalid answer provided to invite", err.Error())
	})
}
