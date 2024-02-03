package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	"github.com/RuanScherer/journey-track-api/adapters/rest/utils"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type ProjectHandler struct {
	projectUseCase usecase.ProjectUseCase
}

func NewProjectHandler(projectUseCase usecase.ProjectUseCase) *ProjectHandler {
	return &ProjectHandler{projectUseCase}
}

func (handler *ProjectHandler) CreateProject(ctx *fiber.Ctx) error {
	req := new(appmodel.CreateProjectRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	err = utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.projectUseCase.CreateProject(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(res)
}

func (handler *ProjectHandler) EditProject(ctx *fiber.Ctx) error {
	req := new(appmodel.EditProjectRequest)
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	req.ActorID = ctx.Locals("sessionUser").(appmodel.AuthUser).ID
	req.ProjectID = ctx.Params("id")

	err = utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.projectUseCase.EditProject(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (handler *ProjectHandler) ShowProject(ctx *fiber.Ctx) error {
	req := &appmodel.ShowProjectRequest{
		ActorID:   ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ProjectID: ctx.Params("id"),
	}

	err := utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.projectUseCase.ShowProject(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (handler *ProjectHandler) GetProjectStats(ctx *fiber.Ctx) error {
	req := &appmodel.GetProjectStatsRequest{
		ActorID:   ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ProjectID: ctx.Params("id"),
	}

	err := utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	res, err := handler.projectUseCase.GetProjectStats(req)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (handler *ProjectHandler) ListProjectsByMember(ctx *fiber.Ctx) error {
	userId := ctx.Locals("sessionUser").(appmodel.AuthUser).ID

	res, err := handler.projectUseCase.ListProjectsByMember(userId)
	if err != nil {
		return err
	}

	// fiber was sending `null` instead of empty array, so I did this
	if len(*res) == 0 {
		return ctx.Status(fiber.StatusOK).JSON([]any{})
	}
	return ctx.Status(fiber.StatusOK).JSON(*res)
}

func (handler *ProjectHandler) DeleteProject(ctx *fiber.Ctx) error {
	req := &appmodel.DeleteProjectRequest{
		ActorID:   ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ProjectID: ctx.Params("id"),
	}

	err := utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	err = handler.projectUseCase.DeleteProject(req)
	if err != nil {
		return err
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (handler *ProjectHandler) InviteProjectMember(ctx *fiber.Ctx) error {
	req := &appmodel.InviteProjectMemberRequest{
		ActorID:   ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ProjectID: ctx.Params("projectId"),
		UserID:    ctx.Params("userId"),
	}

	err := utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	invite, err := handler.projectUseCase.InviteMember(req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusCreated).JSON(invite)
}

func (handler *ProjectHandler) AcceptProjectInvite(ctx *fiber.Ctx) error {
	req := &appmodel.AnswerProjectInviteRequest{}
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}
	req.ProjectID = ctx.Params("projectId")

	err = utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	err = handler.projectUseCase.AcceptInvite(req)
	if err != nil {
		return err
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (handler *ProjectHandler) DeclineProjectInvite(ctx *fiber.Ctx) error {
	req := &appmodel.AnswerProjectInviteRequest{}
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}
	req.ProjectID = ctx.Params("projectId")

	err = utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	err = handler.projectUseCase.DeclineInvite(req)
	if err != nil {
		return err
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}

func (handler *ProjectHandler) RevokeProjectInvite(ctx *fiber.Ctx) error {
	req := &appmodel.RevokeProjectInviteRequest{
		ActorID:         ctx.Locals("sessionUser").(appmodel.AuthUser).ID,
		ProjectInviteID: ctx.Params("id"),
	}

	err := utils.ValidateRequestBody(req)
	if err != nil {
		return err
	}

	err = handler.projectUseCase.RevokeInvite(req)
	if err != nil {
		return err
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
