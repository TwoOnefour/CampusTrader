package jwtUtils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

var JwtKey = []byte(os.Getenv("JWT_KEY"))

func GenerateToken(userID uint, username string) (string, error) {
	claims := AuthClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			Issuer:    "github.com/twoonefour/CampusTrader",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}
