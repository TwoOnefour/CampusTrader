package main

import (
	"CampusTrader/internal/common/database"
	"CampusTrader/internal/common/router"
	"CampusTrader/internal/common/storage"
	"CampusTrader/internal/controller"
	"CampusTrader/internal/service"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	database.InitMySQL()

	r := router.InitRouter(
		controller.NewUserController(service.NewUserService(database.DB)),
		controller.NewProductController(service.NewProductService(database.DB, service.NewLogService(database.DB))),
		controller.NewOrderController(service.NewOrderService(database.DB, service.NewLogService(database.DB))),
		controller.NewImageController(service.NewImageService(storage.NewLocalStorage("./static/uploads", "/static/uploads"))),
		controller.NewStatisticsController(service.NewStatisticsService(database.DB)),
		controller.NewCategoryController(service.NewCategoryService(database.DB)),
	)
	r.Run(":8080")

}
