package helper

import (
	"github.com/golang-jwt/jwt"
	"go-api/model"
	"time"
)

var jwtSecretKey = []byte("jwt-secret-key")
var JWTCookieName = "Signature-Key"

func GenerateJWT(user *model.UserResponse) (string, error) {
	customClaim := model.GojekJWTClaim{
		Roles: []string{"user"},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.TimeFunc().Add(time.Hour * 24).Unix(),
			Id:        user.Id,
			IssuedAt:  jwt.TimeFunc().Unix(),
			Issuer:    "gojekclone",
			NotBefore: jwt.TimeFunc().Unix(),
			Subject:   user.Username,
		},
	}

	payload := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaim)
	token, err := payload.SignedString(jwtSecretKey)
	return token, err
}

func ValidateJWT(tokenString string) (*model.GojekJWTClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &model.GojekJWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	claims, ok := token.Claims.(*model.GojekJWTClaim)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}
