package rest

import (
	"github.com/RuanScherer/journey-track-api/application/rest/handlers"
	"github.com/RuanScherer/journey-track-api/application/rest/middlewares"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/RuanScherer/journey-track-api/infrastructure/db"
	"github.com/RuanScherer/journey-track-api/infrastructure/repository"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App) {
	dbConn := db.GetConnection()

	userRepository := repository.NewUserDBRepository(dbConn)
	projectRepository := repository.NewProjectDBRepository(dbConn)
	projectInviteRepository := repository.NewProjectInviteDBRepository(dbConn)
	eventRepository := repository.NewEventDBRepository(dbConn)

	userUseCase := usecase.NewUserUseCase(userRepository)
	_ = usecase.NewProjectUseCase(
		projectRepository,
		userRepository,
		projectInviteRepository,
		eventRepository,
	)

	userHandler := handlers.NewUserHandler(*userUseCase)

	api := app.Group("api")
	v1 := api.Group("v1")

	v1.Post("/signin", userHandler.SignIn)

	v1.Post("/users/register", userHandler.RegisterUser)
	v1.Patch("/users/:id/verify/:token", userHandler.VerifyUser)

	// auth middleware - separate protected routes
	api.Use(middlewares.HandleAuth)

	// v1.Get("/users/", func(ctx *fiber.Ctx) error {
	// 	return ctx.SendString(ctx.Locals("sessionUser").(model.AuthUser).Name)
	// })
}
