package jwtUtils

import (
	"github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
