package model

import "github.com/golang-jwt/jwt"

type (
	GojekJWTClaim struct {
		Roles []string `json:"roles"`
		jwt.StandardClaims
	}
)
