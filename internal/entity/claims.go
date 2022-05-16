package entity

import "github.com/golang-jwt/jwt"

type Claims struct {
	Email  string `json:"email"`
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}
