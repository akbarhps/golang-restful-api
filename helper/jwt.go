package helper

import (
	"github.com/golang-jwt/jwt"
	"time"
)

var jwtSecretKey = []byte("jwt-secret-key")
var JWTCookieName = "Signature-Key"

func GenerateJWT(uid, username string) (string, error) {
	payload := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ExpiresAt: jwt.TimeFunc().Add(time.Hour * 24).Unix(),
		Id:        uid,
		IssuedAt:  jwt.TimeFunc().Unix(),
		Issuer:    "instapounds",
		NotBefore: jwt.TimeFunc().Unix(),
		Subject:   username,
	})
	token, err := payload.SignedString(jwtSecretKey)
	return token, err
}

func ValidateJWT(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid {
		return nil, err
	}
	return claims, nil
}
