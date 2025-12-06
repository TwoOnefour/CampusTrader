package router

import (
	"CampusTrader/internal/controller"
	"CampusTrader/internal/middleware/auth"
	"github.com/gin-gonic/gin"
)

func InitRouter(userCtrl *controller.UserController) *gin.Engine {
	r := gin.Default()

	r.Use(Cors())
	// 4. 路由分组 (推荐)
	// 这里的 /api/v1 是所有接口的前缀
	apiGroup := r.Group("/api/v1")
	{
		// Public 路由 (不需要登录)
		apiGroup.POST("/login", userCtrl.Login)
		apiGroup.POST("/register", userCtrl.Register)

		// Private 路由 (需要 JWT 验证)
		// 使用 Group 嵌套中间件
		authGroup := apiGroup.Group("/")
		authGroup.Use(auth.JWTAuthMiddleware())
		{
			authGroup.GET("/me", userCtrl.Me)
		}
	}

	return r
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
