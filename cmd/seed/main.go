package main

import (
	"CampusTrader/internal/common/database"
	"CampusTrader/internal/model"
	"fmt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
)

func main() {
	// 1. 加载环境变量 (为了获取 DSN)
	if err := godotenv.Load(); err != nil {
		log.Println("注意: 没有找到 .env 文件，尝试直接读取环境变量")
	}

	// 2. 初始化数据库连接
	database.InitMySQL()

	db := database.DB

	fmt.Println("🌱 开始播种数据...")

	// 3. 清理旧数据 (可选，防止重复运行报错)
	cleanData(db)

	// 4. 创建测试用户
	// 注意：必须手动加密密码，因为直接插入数据库不会经过 Service 层
	password := "password123"
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := []model.User{
		{
			Username: "testuser",
			Password: string(hashedPwd),
			Nickname: "测试用户",
			Email:    "test@example.com",
			Phone:    "13800138000",
		},
		{
			Username: "testuser2",
			Password: string(hashedPwd),
			Nickname: "测试用户2",
			Email:    "test2@example.com",
			Phone:    "13800138000",
		},
	}
	if err := db.Create(&user).Error; err != nil {
		panic(err)
	}
	// 5. 创建分类
	category := model.Category{
		Name: "电子数码",
	}
	if err := db.Create(&category).Error; err != nil {
		panic(err)
	}
	fmt.Printf("✅ 分类创建成功: %s (ID: %d)\n", category.Name, category.Id)

	// 6. 创建商品
	products := []model.Product{
		{
			Name:        "MacBook Pro M3",
			Description: "几乎全新，仅循环充电 10 次，箱说全。",
			Price:       12999.00,
			CategoryId:  category.Id,
			SellerId:    user[0].ID, // 关联上面创建的用户
			Status:      "available",
			Condition:   "like_new",
			ImageUrl:    "https://images.unsplash.com/photo-1724859234679-964acf07b126?w=500&auto=format&fit=crop&q=60&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxzZWFyY2h8Mnx8TWFjQm9vayUyMFBybyUyME0zfGVufDB8fDB8fHww",
		},
		{
			Name:        "Sony WH-1000XM5",
			Description: "降噪耳机，音质无敌，考研党必备。",
			Price:       1899.00,
			CategoryId:  category.Id,
			SellerId:    user[0].ID,
			Status:      "available",
			Condition:   "good",
			ImageUrl:    "https://images.unsplash.com/photo-1618366712010-f4ae9c647dcb?ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60",
		},
		{
			Name:        "IKEA 台灯",
			Description: "毕业带不走，低价出。",
			Price:       25.00,
			CategoryId:  category.Id,
			SellerId:    user[0].ID,
			Status:      "sold", // 这个已售出，测试前端是否变灰
			Condition:   "fair",
			ImageUrl:    "https://images.unsplash.com/photo-1705081820804-22e5a09149d4?q=80&w=1935&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D",
		},
	}

	if err := db.Create(&products).Error; err != nil {
		panic(err)
	}
	fmt.Printf("✅ 成功插入 %d 个商品\n", len(products))

	fmt.Println("🎉 数据播种完成！现在可以启动后端并刷新前端页面了。")
}

func cleanData(db *gorm.DB) {
	// 硬删除清空表，根据你的表名调整
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM categories")
	db.Exec("DELETE FROM users")
	// 重置自增 ID (MySQL)
	db.Exec("ALTER TABLE products AUTO_INCREMENT = 1")
	db.Exec("ALTER TABLE categories AUTO_INCREMENT = 1")
	db.Exec("ALTER TABLE users AUTO_INCREMENT = 1")
}

func dropData(db *gorm.DB) {
	db.Exec("DROP TABLE products")
	db.Exec("DROP TABLE categories")
	db.Exec("DROP TABLE users")
}
