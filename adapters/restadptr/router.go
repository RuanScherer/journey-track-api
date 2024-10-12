package restadptr

import (
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/handler"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("api")
	v1 := api.Group("v1")

	v1.Post("/signin", handler.NewSignInHandler().Handle)
	v1.Post("/signout", handler.NewSignOutHandler().Handle)

	v1.Post("/users/request-password-reset", handler.NewRequestUserPasswordResetHandler().Handle)
	v1.Patch("/users/:id/reset-password/:token", handler.NewResetUserPassword().Handle)
	v1.Post("/users/register", handler.NewRegisterUserHandler().Handle)
	v1.Patch("/users/:id/verify/:token", handler.NewVerifyUserHandler().Handle)

	v1.Get("/projects/:projectId/invites/:token", handler.NewShowInvitationByProjectAndTokenHandler().Handle)

	// auth middleware - separate protected routes
	api.Use(middleware.HandleAuth)

	v1.Put("/users/edit-profile", handler.NewEditUserHandler().Handle)
	v1.Get("/users/profile", handler.NewShowUserHandler().Handle)
	v1.Get("/users/search", handler.NewSearchUsersHandler().Handle)

	v1.Post("/projects/create", handler.NewCreateProjectHandler().Handle)
	v1.Put("/projects/:id/edit", handler.NewEditProjectHandler().Handle)
	v1.Get("/projects/:id", handler.NewShowProjectHandler().Handle)
	v1.Get("/projects/:id/stats", handler.NewGetProjectStatsHandler().Handle)
	v1.Get("/projects", handler.NewListProjectsByMemberHandler().Handle)
	v1.Delete("/projects/:id", handler.NewDeleteProjectHandler().Handle)

	v1.Get("/projects/:projectId/invites", handler.NewListProjectInvitesHandler().Handle)
	v1.Post("/projects/:projectId/invite", handler.NewInviteProjectMembersHandler().Handle)
	v1.Patch("/projects/:projectId/invites/accept", handler.NewAcceptProjectInviteHandler().Handle)
	v1.Patch("/projects/:projectId/invites/decline", handler.NewDeclineProjectInviteHandler().Handle)
	v1.Delete("/projects/invites/:id/revoke", handler.NewRevokeProjectInviteHandler().Handle)
}
