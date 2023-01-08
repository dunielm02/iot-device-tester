package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const SecretKey = "2L1DL3AjCfPjRcFh7L5JFRV7VijapVWnr13ReVsBNiY2L1DL3AjCfPVWnr13ReVs"

func GenerateJwt(username string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	})

	return claims.SignedString([]byte(SecretKey))
}

func GetIdFromJwt(cookie string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	tokenWithClaims := token.Claims.(*jwt.RegisteredClaims)

	return tokenWithClaims.Issuer, nil
}
