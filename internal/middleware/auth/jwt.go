package auth

import (
	"CampusTrader/internal/common/response"
	"CampusTrader/internal/util/jwtUtils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

var jwtKey = jwtUtils.JwtKey

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(ctx, http.StatusUnauthorized, "无效的token")
			ctx.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(ctx, http.StatusUnauthorized, "无效的token")
			ctx.Abort()
			return
		}
		tokenString := parts[1]
		
		token, err := jwt.ParseWithClaims(tokenString, &jwtUtils.AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			response.Error(ctx, http.StatusUnauthorized, "无效的token")
			ctx.Abort()
			return
		}
		if claims, ok := token.Claims.(*jwtUtils.AuthClaims); ok && token.Valid {
			ctx.Set("userID", claims.UserID)
			ctx.Set("username", claims.Username)
		}
	}
}
