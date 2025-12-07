package router

import (
	"CampusTrader/internal/controller"
	"CampusTrader/internal/middleware/auth"
	"github.com/gin-gonic/gin"
)

func InitRouter(
	userCtrl *controller.UserController,
	productCtrl *controller.ProductController, // 新增
	orderCtrl *controller.OrderController,     // 新增
	imageCtrl *controller.ImageController,
) *gin.Engine {
	r := gin.Default()
	r.Use(Cors())

	apiGroup := r.Group("/api/v1")
	{
		// User
		apiGroup.POST("/login", userCtrl.Login)
		apiGroup.POST("/register", userCtrl.Register)

		// Product (无需登录)
		apiGroup.GET("/product/list", productCtrl.ListProducts)
		apiGroup.GET("/category/list")
		// Private Group
		authGroup := apiGroup.Group("/")
		authGroup.Use(auth.JWTAuthMiddleware())
		{
			authGroup.GET("/me", userCtrl.Me)

			// Order (需要登录)
			authGroup.POST("/order/create", orderCtrl.Order)
			authGroup.GET("/order/my", orderCtrl.ListOrder)

			// Product (发布需要登录)

			authGroup.POST("/product/create", productCtrl.CreateProduct)
			// image上传要登陆
			authGroup.POST("/image/upload", imageCtrl.Upload)
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
