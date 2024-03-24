package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/postgres"
	"github.com/RuanScherer/journey-track-api/adapters/postgres/repository"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type ListProjectsByMemberHandler struct {
	useCase usecase.ListProjectsByMemberUseCase
}

func NewListProjectsByMemberHandler() *ListProjectsByMemberHandler {
	projectRepository := repository.NewProjectPostgresRepository(postgres.GetConnection())
	useCase := *usecase.NewListProjectsByMemberUseCase(projectRepository)
	return &ListProjectsByMemberHandler{useCase: useCase}
}

func (handler *ListProjectsByMemberHandler) Handle(ctx *fiber.Ctx) error {
	userId := ctx.Locals("sessionUser").(appmodel.AuthUser).ID

	res, err := handler.useCase.Execute(userId)
	if err != nil {
		return err
	}

	// fiber was sending `null` instead of empty array, so I did this
	if len(*res) == 0 {
		return ctx.Status(fiber.StatusOK).JSON([]any{})
	}
	return ctx.Status(fiber.StatusOK).JSON(*res)
}
