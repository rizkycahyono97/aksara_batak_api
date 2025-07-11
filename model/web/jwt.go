package web

import "github.com/golang-jwt/jwt/v5"

type JwtCustomClaims struct {
	UUID string `json:"uuid"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}
