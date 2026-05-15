package model

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TokenExp = time.Hour * 24
const CookieName = "auth_token"

type Claims struct {
	jwt.RegisteredClaims
	UserID int `json:"user_id"`
}

func BuildJWTString(userID int, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GetUserID(tokenString, secretKey string) (int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return -1, err
	}

	if !token.Valid {
		return -1, fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}

func GenerateUserID() int {
	return int(time.Now().UnixNano())
}
