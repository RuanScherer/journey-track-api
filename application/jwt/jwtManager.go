package jwt

import (
	"errors"
	"time"

	appmodel "github.com/RuanScherer/journey-track-api/application/model"
	"github.com/RuanScherer/journey-track-api/config"
	"github.com/RuanScherer/journey-track-api/domain/model"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ExpirationTime = time.Hour * 24 * 7
)

type Manager interface {
	CreateJwtFromUser(user *model.User) (string, error)
	GetJwtClaims(token string) (*appmodel.JwtClaims, error)
}

type DefaultManager struct{}

func NewDefaultManager() *DefaultManager {
	return &DefaultManager{}
}

func (manager *DefaultManager) CreateJwtFromUser(user *model.User) (string, error) {
	jwtClaims := appmodel.JwtClaims{
		User: appmodel.AuthUser{
			ID:    user.ID,
			Email: *user.Email,
			Name:  user.Name,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ExpirationTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	appConfig := config.GetAppConfig()
	jwtString, err := token.SignedString([]byte(appConfig.JwtSecret))
	if err != nil {
		return "", errors.New("error creating access token")
	}

	return jwtString, nil
}

func (manager *DefaultManager) GetJwtClaims(token string) (*appmodel.JwtClaims, error) {
	if token == "" {
		return nil, appmodel.NewAppError("missing_access_token", "missing access token", appmodel.ErrorTypeAuthentication)
	}

	parsedToken, err := jwt.ParseWithClaims(token, &appmodel.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetAppConfig().JwtSecret), nil
	})
	if err != nil {
		return nil, appmodel.NewAppError("invalid_access_token", "invalid access token", appmodel.ErrorTypeAuthentication)
	}

	claims, ok := parsedToken.Claims.(*appmodel.JwtClaims)
	if !ok {
		return nil, appmodel.NewAppError("invalid_access_token", "invalid access token", appmodel.ErrorTypeAuthentication)
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, appmodel.NewAppError("expired_access_token", "expired access token", appmodel.ErrorTypeAuthentication)
	}

	return claims, nil
}
