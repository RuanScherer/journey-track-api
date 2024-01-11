package middlewares

import (
	"time"

	"github.com/RuanScherer/journey-track-api/application/utils"
	"github.com/gofiber/fiber/v2"
)

func HandleAuth(ctx *fiber.Ctx) error {
	claims, err := utils.GetJwtClaims(ctx.Cookies("access_token"))
	if err != nil {
		return err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return fiber.NewError(fiber.StatusUnauthorized, "expired access token")
	}

	ctx.Locals("sessionUser", claims.User)
	return ctx.Next()
}
