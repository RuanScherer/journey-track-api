package middleware

import (
	"github.com/RuanScherer/journey-track-api/adapters/rest/model"
	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/gofiber/fiber/v2"
)

func HandleError(ctx *fiber.Ctx, err error) error {
	if err, ok := err.(*model.RestApiError); ok {
		return ctx.Status(err.StatusCode).JSON(err)
	}

	if err, ok := err.(*appmodel.AppError); ok {
		statusCode := getStatusCodeFromAppError(err)
		return ctx.Status(statusCode).JSON(model.NewRestApiError(statusCode, err))
	}

	if err, ok := err.(*fiber.Error); ok {
		appError := &appmodel.AppError{
			Code:    "unexpected_error",
			Message: err.Message,
		}
		return ctx.Status(err.Code).JSON(model.NewRestApiError(err.Code, appError))
	}

	appError := &appmodel.AppError{
		Code:    "unexpected_error",
		Message: err.Error(),
	}
	statusCode := fiber.StatusInternalServerError
	return ctx.Status(statusCode).JSON(model.NewRestApiError(statusCode, appError))
}

func getStatusCodeFromAppError(err *appmodel.AppError) int {
	switch err.Type {
	case appmodel.ErrorTypeValidation, appmodel.ErrorTypeRequest:
		return fiber.StatusBadRequest
	case appmodel.ErrorTypeDatabase, appmodel.ErrorTypeServer:
		return fiber.StatusInternalServerError
	case appmodel.ErrorTypeAuthentication:
		return fiber.StatusUnauthorized
	default:
		return fiber.StatusInternalServerError
	}
}
