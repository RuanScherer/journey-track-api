package utils

import (
	"errors"
	"time"

	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/config"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJwtFromUser(user *model.User) (string, error) {
	jwtClaims := appmodel.JwtClaims{
		User: appmodel.AuthUser{
			ID:    user.ID,
			Email: *user.Email,
			Name:  user.Name,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	appConfig := config.GetAppConfig()
	jwtString, err := jwt.SignedString([]byte(appConfig.JwtSecret))
	if err != nil {
		return "", errors.New("error creating access token")
	}

	return jwtString, nil
}
