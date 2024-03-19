package rest

import (
	"github.com/RuanScherer/journey-track-api/adapters/rest/handlers"
	"github.com/RuanScherer/journey-track-api/adapters/rest/middlewares"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	api := app.Group("api")
	v1 := api.Group("v1")

	v1.Post("/signin", handlers.NewSignInHandler().Handle)
	v1.Post("/signout", handlers.NewSignOutHandler().Handle)

	v1.Post("/users/request-password-reset", handlers.NewRequestUserPasswordResetHandler().Handle)
	v1.Patch("/users/:id/reset-password/:token", handlers.NewResetUserPassword().Handle)
	v1.Post("/users/register", handlers.NewRegisterUserHandler().Handle)
	v1.Patch("/users/:id/verify/:token", handlers.NewVerifyUserHandler().Handle)

	// auth middleware - separate protected routes
	api.Use(middlewares.HandleAuth)

	v1.Put("/users/edit-profile", handlers.NewEditUserHandler().Handle)
	v1.Get("/users/profile", handlers.NewShowUserHandler().Handle)
	v1.Get("/users/search", handlers.NewSearchUsersHandler().Handle)

	v1.Post("/projects/create", handlers.NewCreateProjectHandler().Handle)
	v1.Put("/projects/:id/edit", handlers.NewEditProjectHandler().Handle)
	v1.Get("/projects/:id", handlers.NewShowProjectHandler().Handle)
	v1.Get("/projects/:id/stats", handlers.NewGetProjectStatsHandler().Handle)
	v1.Get("/projects", handlers.NewListProjectsByMemberHandler().Handle)
	v1.Delete("/projects/:id", handlers.NewDeleteProjectHandler().Handle)

	v1.Get("/projects/:projectId/invites", handlers.NewListProjectInvitesHandler().Handle)
	v1.Post("/projects/:projectId/invite", handlers.NewInviteProjectMembersHandler().Handle)
	v1.Patch("/projects/:projectId/invites/accept", handlers.NewAcceptProjectInviteHandler().Handle)
	v1.Patch("/projects/:projectId/invites/decline", handlers.NewDeclineProjectInviteHandler().Handle)
	v1.Delete("/projects/invites/:id/revoke", handlers.NewRevokeProjectInviteHandler().Handle)
}
