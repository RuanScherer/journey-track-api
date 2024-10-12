package middleware

import (
	"time"

	"github.com/RuanScherer/journey-track-api/application/jwt"
	"github.com/gofiber/fiber/v2"
)

func ExpireAccessTokenCookie(ctx *fiber.Ctx) {
	// needed to expire the cookie this way due to a bug in fiber
	// https://github.com/gofiber/fiber/issues/1127
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "deleted",
		HTTPOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(-3 * time.Second),
	})
}

func HandleAuth(ctx *fiber.Ctx) error {
	claims, err := jwt.NewDefaultManager().GetJwtClaims(ctx.Cookies("access_token"))
	if err != nil {
		ExpireAccessTokenCookie(ctx)
		return err
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		ExpireAccessTokenCookie(ctx)
		return fiber.NewError(fiber.StatusUnauthorized, "expired access token")
	}

	ctx.Locals("sessionUser", claims.User)
	return ctx.Next()
}
