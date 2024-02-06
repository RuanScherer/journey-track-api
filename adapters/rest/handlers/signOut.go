package handlers

import (
	"github.com/RuanScherer/journey-track-api/adapters/rest/middlewares"
	"github.com/gofiber/fiber/v2"
)

type SignOutHandler struct{}

func NewSignOutHandler() *SignOutHandler {
	return &SignOutHandler{}
}

func (handler *SignOutHandler) Handle(ctx *fiber.Ctx) error {
	middlewares.ExpireAccessTokenCookie(ctx)
	ctx.Status(fiber.StatusNoContent)
	return nil
}
