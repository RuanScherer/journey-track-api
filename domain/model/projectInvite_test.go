package model

import "testing"

func TestNewProjectInvite(t *testing.T) {
	t.Run("should get error when provided project is invalid", func(t *testing.T) {
		_, err := NewProjectInvite(&Project{}, &User{})
		if err == nil {
			t.Error("should get error when provided project is invalid")
		}
	})

	t.Run("should get error when provided user is invalid", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		project, _ := NewProject("test", projectOwner)
		_, err := NewProjectInvite(project, &User{})

		if err == nil {
			t.Error("should get error when provided user is invalid")
		}
	})

	t.Run("should get error when user is already a member", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, err := NewProjectInvite(project, userToInvite)
		if err != nil {
			t.Error("should return project invite when provided project and user are valid")
		}

		invite.Accept(invite.Token)
		project.Members = append(project.Members, userToInvite)

		_, err = NewProjectInvite(project, userToInvite)
		if err == nil || err.Error() != "user is already a member of the project" {
			t.Error("should get error when user is already a member")
		}
	})

	t.Run("should get error when user already has a pending invite for project", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, err := NewProjectInvite(project, userToInvite)
		if err != nil {
			t.Error("should return project invite when provided project and user are valid")
		}
		userToInvite.ProjectInvites = append(userToInvite.ProjectInvites, invite)

		_, err = NewProjectInvite(project, userToInvite)
		if err == nil || err.Error() != "user already has a pending invite for the project" {
			t.Error("should get error when user already has a pending invite for project")
		}
	})

	t.Run("should return project invite when provided project and user are valid", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		_, err := NewProjectInvite(project, userToInvite)

		if err != nil {
			t.Error("should return project invite when provided project and user are valid")
		}
	})

	t.Run("should return project invite with pending status", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		if invite.Status != ProjectInviteStatusPending {
			t.Error("should return project invite with pending status")
		}
	})

	t.Run("should return project invite with token", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		if invite.Token == "" {
			t.Error("should return project invite with token")
		}
	})
}

func TestAccept(t *testing.T) {
	t.Run("should get error when provided token is invalid", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		err := invite.Accept("invalid-token")
		if err == nil || err.Error() != "invalid token provided to answer invite" {
			t.Error("should get error when provided token is invalid")
		}
	})

	t.Run("should get error when invited is not pending", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)
		invite.Status = ProjectInviteStatusAccepted

		err := invite.Accept(invite.Token)
		if err == nil || err.Error() != "invite already answered or revoked" {
			t.Error("should get error when invited is not pending")
		}
	})

	t.Run("should accept invite", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		err := invite.Accept(invite.Token)
		if err != nil {
			t.Error("should accept invite")
		}

		if invite.Status != ProjectInviteStatusAccepted {
			t.Error("should have accepted status after accepting invite")
		}
	})
}

func TestDecline(t *testing.T) {
	t.Run("should get error when provided token is invalid", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		err := invite.Decline("invalid-token")
		if err == nil || err.Error() != "invalid token provided to answer invite" {
			t.Error("should get error when provided token is invalid")
		}
	})

	t.Run("should get error when invited is not pending", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)
		invite.Status = ProjectInviteStatusAccepted

		err := invite.Decline(invite.Token)
		if err == nil || err.Error() != "invite already answered or revoked" {
			t.Error("should get error when invited is not pending")
		}
	})

	t.Run("should decline invite", func(t *testing.T) {
		projectOwner, _ := NewUser("owner@example.com", "Owner", "pass1234")
		projectOwner.Verify(projectOwner.VerificationToken)
		project, _ := NewProject("test", projectOwner)

		userToInvite, _ := NewUser("member@example.com", "Member", "pass4321")
		invite, _ := NewProjectInvite(project, userToInvite)

		err := invite.Decline(invite.Token)
		if err != nil {
			t.Error("should decline invite")
		}

		if invite.Status != ProjectInviteStatusDeclined {
			t.Error("should have declined status after accepting invite")
		}
	})
}
