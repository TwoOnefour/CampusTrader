package main

import (
	"CampusTrader/internal/common/database"
	"CampusTrader/internal/common/router"
	"CampusTrader/internal/common/storage"
	"CampusTrader/internal/controller"
	"CampusTrader/internal/service"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	database.InitMySQL()

	logService := service.NewLogService(database.DB)
	userService := service.NewUserService(database.DB)

	orderService := service.NewOrderService(database.DB, logService)
	storageService := storage.NewLocalStorage("./static/uploads", "/static/uploads")
	imageService := service.NewImageService(storageService)
	categoryService := service.NewCategoryService(database.DB)
	statService := service.NewStatisticsService(database.DB)
	productService := service.NewProductService(database.DB, logService, statService)
	r := router.InitRouter(
		controller.NewUserController(userService),
		controller.NewProductController(productService, statService),
		controller.NewOrderController(orderService),
		controller.NewImageController(imageService),
		controller.NewStatisticsController(statService),
		controller.NewCategoryController(categoryService),
	)
	r.Run(":8080")
}
