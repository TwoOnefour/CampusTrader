package router

import (
	"CampusTrader/internal/controller"
	"CampusTrader/internal/middleware/auth"
	"github.com/gin-gonic/gin"
)

func InitRouter(
	userCtrl *controller.UserController,
	productCtrl *controller.ProductController, // 新增
	orderCtrl *controller.OrderController, // 新增
	imageCtrl *controller.ImageController,
) *gin.Engine {
	r := gin.Default()
	r.Use(Cors())

	apiGroup := r.Group("/api/v1")
	{
		// User
		apiGroup.POST("/login", userCtrl.Login)
		apiGroup.POST("/register", userCtrl.Register)

		apiGroup.GET("/category/list")

		// Product (无需登录)
		{
			apiGroup.GET("/product/search", productCtrl.SearchProducts)
			apiGroup.GET("/product/suggestion", productCtrl.SearchProductsSuggestion)
			apiGroup.GET("/product/list", productCtrl.ListProducts)
		}

		// Private Group
		authGroup := apiGroup.Group("/")
		authGroup.Use(auth.JWTAuthMiddleware())
		{
			authGroup.GET("/me", userCtrl.Me)

			// Order (需要登录)
			authGroup.POST("/order/create", orderCtrl.Order)
			authGroup.GET("/order/my", orderCtrl.ListOrder)
			// image上传要登陆
			authGroup.POST("/image/upload", imageCtrl.Upload)

			authGroup.GET("/product/my", productCtrl.ListMyProducts)
			authGroup.POST("/product/create", productCtrl.CreateProduct)
			authGroup.POST("/product/drop", productCtrl.DropProduct)
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
