package handlers

import (
	"strconv"
	"strings"

	"github.com/RuanScherer/journey-track-api/adapters/db"
	"github.com/RuanScherer/journey-track-api/adapters/db/repositories"
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	"github.com/RuanScherer/journey-track-api/adapters/rest/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type SearchUsersHandler struct {
	useCase usecase.SearchUsersUseCase
}

func NewSearchUsersHandler() *SearchUsersHandler {
	userRepository := repositories.NewUserDBRepository(db.GetConnection())
	useCase := *usecase.NewSearchUsersUseCase(userRepository)
	return &SearchUsersHandler{useCase: useCase}
}

func (handler *SearchUsersHandler) Handle(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size", "10"))
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	req := &appmodel.SearchUsersRequest{
		ActorID:            ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ExcludedProjectIDs: strings.Split(ctx.Query("excluded_project_ids"), ","),
		Email:              ctx.Query("email"),
		Page:               page,
		PageSize:           pageSize,
	}

	err = utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.useCase.Execute(req)
	if err != nil {
		return err
	}

	// fiber was sending `null` instead of empty array, so I did this
	if len(*res) == 0 {
		return ctx.JSON([]any{})
	}
	return ctx.JSON(*res)
}
