package model

import "github.com/golang-jwt/jwt/v5"

type JwtClaims struct {
	User AuthUser `json:"user"`
	jwt.RegisteredClaims
}

type AuthUser struct {
	ID    string `json:"id"`
	Email string `json:"smtpemail"`
	Name  string `json:"name"`
}
