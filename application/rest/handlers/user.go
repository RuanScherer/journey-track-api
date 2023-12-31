package handlers

import (
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/application/rest/model"
	"github.com/RuanScherer/journey-track-api/application/usecase"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase}
}

func (handler *UserHandler) RegisterUser(ctx *fiber.Ctx) error {
	registerUserRequest := &appmodel.RegisterUserRequest{}
	err := ctx.BodyParser(registerUserRequest)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	registerUserResponse, appErr := handler.userUseCase.RegisterUser(registerUserRequest)
	if appErr != nil {
		return appErr
	}

	ctx.Status(fiber.StatusCreated)
	return ctx.JSON(registerUserResponse)
}

func (handler *UserHandler) VerifyUser(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	token := ctx.Params("token")

	err := handler.userUseCase.VerifyUser(&appmodel.VerifyUserRequest{
		UserID:            userId,
		VerificationToken: token,
	})
	if err != nil {
		return err
	}

	ctx.Status(fiber.StatusOK)
	return nil
}

func (handler *UserHandler) SignIn(ctx *fiber.Ctx) error {
	signInRequest := &appmodel.SignInRequest{}
	err := ctx.BodyParser(signInRequest)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	signInResponse, appErr := handler.userUseCase.SignIn(signInRequest)
	if appErr != nil {
		return appErr
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    signInResponse.AccessToken,
		HTTPOnly: true,
		Path:     "/",
	})
	return ctx.JSON(signInResponse)
}

func (handler *UserHandler) EditUser(ctx *fiber.Ctx) error {
	editUserRequest := &appmodel.EditUserRequest{}
	err := ctx.BodyParser(editUserRequest)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	editUserRequest.UserID = ctx.Locals("sessionUser").(appmodel.AuthUser).ID
	response, appErr := handler.userUseCase.EditUser(editUserRequest)
	if appErr != nil {
		return appErr
	}
	return ctx.JSON(response)
}

func (handler *UserHandler) RequestPasswordReset(ctx *fiber.Ctx) error {
	req := &appmodel.RequestPasswordResetRequest{}
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	err = handler.userUseCase.RequestUserPasswordReset(req)
	return err
}

func (handler *UserHandler) ResetPassword(ctx *fiber.Ctx) error {
	req := &appmodel.PasswordResetRequest{}
	err := ctx.BodyParser(req)
	if err != nil {
		return model.NewRestApiError(fiber.StatusBadRequest, appmodel.ErrInvalidReqData)
	}

	req.UserID = ctx.Params("id")
	req.PasswordResetToken = ctx.Params("token")

	err = handler.userUseCase.ResetUserPassword(req)
	return err
}

func (handler *UserHandler) ShowUser(ctx *fiber.Ctx) error {
	userID := ctx.Locals("sessionUser").(appmodel.AuthUser).ID

	response, err := handler.userUseCase.ShowUser(userID)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}
