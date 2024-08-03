package api

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var secretKey string = "secret"

func getRole(userId string) string {
	if userId == "1" {
		return "admin"
	} else {
		return "user"
	}
}

func GenerateToken(userId string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"iss": "mqchat",
		"aud": getRole(userId),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})

	fmt.Printf("Token claims added: %+v\n", claims)

	token, _ := claims.SignedString([]byte(secretKey))
	return token, nil
}

func ValidateToken(token string) (bool, error) {
	claims := jwt.MapClaims{}
	t, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return false, err
	}

	if !t.Valid {
		return false, nil
	}

	return true, nil
}



