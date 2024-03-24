package handler

import (
	"github.com/RuanScherer/journey-track-api/adapters/rest/middleware"
	"github.com/gofiber/fiber/v2"
)

type SignOutHandler struct{}

func NewSignOutHandler() *SignOutHandler {
	return &SignOutHandler{}
}

func (handler *SignOutHandler) Handle(ctx *fiber.Ctx) error {
	middleware.ExpireAccessTokenCookie(ctx)
	ctx.Status(fiber.StatusNoContent)
	return nil
}
