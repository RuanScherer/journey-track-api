package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr"
	"github.com/RuanScherer/journey-track-api/adapters/postgresadptr/repository"
	"github.com/RuanScherer/journey-track-api/adapters/restadptr/validator"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/gofiber/fiber/v2"
)

type ListProjectInvitesHandler struct {
	useCase *usecase.ListProjectInvitesUseCase
}

func NewListProjectInvitesHandler() *ListProjectInvitesHandler {
	db := postgresadptr.GetConnection()
	projectInviteRepository := repository.NewProjectInvitePostgresRepository(db)
	projectRepository := repository.NewProjectPostgresRepository(db)
	useCase := usecase.NewListProjectInvitesUseCase(projectInviteRepository, projectRepository)
	return &ListProjectInvitesHandler{useCase}
}

func (handler *ListProjectInvitesHandler) Handle(ctx *fiber.Ctx) error {
	req := &appmodel.ListProjectInvitesRequest{
		ActorID:   ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ProjectID: ctx.Params("projectId"),
		Status:    ctx.Query("status", model.ProjectInviteStatusPending),
	}

	err := validator.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}

	// fiber was sending `null` instead of empty array, so I did this
	if len(*res) == 0 {
		return ctx.Status(fiber.StatusOK).JSON([]any{})
	}
	return ctx.Status(fiber.StatusOK).JSON(*res)
}
