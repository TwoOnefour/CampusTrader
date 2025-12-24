package router

import (
	"CampusTrader/internal/assets"
	"CampusTrader/internal/controller"
	"CampusTrader/internal/middleware/auth"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitRouter(
	userCtrl *controller.UserController,
	productCtrl *controller.ProductController,
	orderCtrl *controller.OrderController,
	imageCtrl *controller.ImageController,
	statisticsCtrl *controller.StatisticsController,
	categoryCtrl *controller.CategoryController,
) *gin.Engine {
	r := gin.Default()
	r.Use(Cors()) // 跨域中间件

	// 建议统一使用 v1 版本组
	v1 := r.Group("/api/v1")
	{
		authGroup := v1.Group("/auth")
		{
			// 虽然是 UserCtrl 处理，但语义上是 Auth
			authGroup.POST("/login", userCtrl.Login)
			authGroup.POST("/register", userCtrl.Register)
		}

		userPublicGroup := v1.Group("/users")
		{
			userPublicGroup.GET("/:id/rating", statisticsCtrl.GetUserRating)
			userPublicGroup.GET("/:id/trade-stats", statisticsCtrl.GetUserCompletedOrderRecord)
		}

		productGroup := v1.Group("/products")
		{
			productGroup.GET("", productCtrl.ListProducts)
			productGroup.GET("/search", productCtrl.SearchProducts)
			productGroup.GET("/suggestion", productCtrl.SearchProductsSuggestion)
		}

		categoryGroup := v1.Group("/categories")
		{
			categoryGroup.GET("", categoryCtrl.ListCategory)               // 原来的 /category/list 改为 GET /categories
			categoryGroup.GET("/popular", statisticsCtrl.GetHotCategories) // 热门分类
		}

		privateGroup := v1.Group("/")
		privateGroup.Use(auth.JWTAuthMiddleware())
		{

			privateGroup.GET("/users/me", userCtrl.Me)
			privateGroup.GET("/users/me/products", productCtrl.ListMyProducts)
			privateGroup.POST("/products", productCtrl.CreateProduct)    // 创建商品用 POST /products
			privateGroup.POST("/products/drop", productCtrl.DropProduct) // 或者 DELETE /products/:id

			orderGroup := privateGroup.Group("/orders")
			{
				orderGroup.POST("", orderCtrl.Order)
				orderGroup.GET("/my", orderCtrl.ListOrder)
			}

			privateGroup.POST("/images", imageCtrl.Upload)
		}
	}
	staticFS := assets.GetFileSystem()
	fileServer := http.FileServer(http.FS(staticFS))
	r.Static("/static", "./static")
	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") || strings.HasPrefix(c.Request.URL.Path, "/static") {
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}
		c.JSON(404, gin.H{"code": 404, "msg": "Not Found"})
	})
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
