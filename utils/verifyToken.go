package utils

import (
	"githb.com/demo-employee-api/internal/entity"
	"github.com/golang-jwt/jwt"
)

func ExtractToken(tokenStr string, JwtKey string) (*entity.Claims, error) {
	claims := &entity.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(JwtKey), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return &entity.Claims{}, err
		}
		return &entity.Claims{}, err
	}

	if !tkn.Valid {
		return &entity.Claims{}, err
	}

	return claims, nil
}
