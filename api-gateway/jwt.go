package main 

import (

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjAxODc0OTgsInVzZXJfaWQiOjEyfQ.hcrK9ay082UFc4FoMxVQN_AqkPg_gkl9vvoH7nSzpeE")

type Claims struct {
	jwt.StandardClaims
	UserID uint `json:"user_id"`
}

func ValidateToken(tokenStr string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.NewValidationError("invalid token", jwt.ValidationErrorSignatureInvalid)
	}
	return claims, nil
}
