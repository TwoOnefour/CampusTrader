package main

import (
	"CampusTrader/internal/common/database"
	"CampusTrader/internal/model"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// 1. 加载环境变量
	if err := godotenv.Load(); err != nil {
		log.Println("注意: 没有找到 .env 文件，尝试直接读取环境变量")
	}

	// 2. 初始化数据库
	database.InitMySQL()
	db := database.DB

	fmt.Println("🌱 开始播种数据...")

	// 3. 清理旧数据 (顺序很重要，先删子表)
	//cleanData(db)

	// 4. 创建测试用户
	password := "password123"
	hashedPwd, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	users := []model.User{
		{
			Username: "testuser",
			Password: string(hashedPwd),
			Nickname: "及格万岁",
			Email:    "test@example.com",
			Phone:    "13800138001",
		},
		{
			Username: "testuser2",
			Password: string(hashedPwd),
			Nickname: "富婆通讯录",
			Email:    "test2@example.com",
			Phone:    "13800138002",
		},
	}
	if err := db.Create(&users).Error; err != nil {
		log.Println(err)
	}

	// 5. 创建分类
	category := model.Category{Name: "电子数码"}
	categoryBook := model.Category{Name: "图书教材"}
	if err := db.Create(&category).Error; err != nil {
		log.Println(err)
	}
	if err := db.Create(&categoryBook).Error; err != nil {
		log.Println(err)
	}

	// 6. 创建商品
	products := []model.Product{
		{
			Name:        "MacBook Pro M3",
			Description: "几乎全新，仅循环充电 10 次，箱说全。",
			Price:       12999.00,
			CategoryId:  category.Id,
			SellerId:    users[0].ID,
			Status:      "available",
			Condition:   "like_new",
			ImageUrl:    "https://images.unsplash.com/photo-1517336714731-489689fd1ca4?w=500&auto=format&fit=crop&q=60",
		},
		{
			Name:        "Sony WH-1000XM5",
			Description: "降噪耳机，音质无敌，考研党必备。",
			Price:       1899.00,
			CategoryId:  category.Id,
			SellerId:    users[0].ID,
			Status:      "sold", // 已售出
			Condition:   "good",
			ImageUrl:    "https://images.unsplash.com/photo-1618366712010-f4ae9c647dcb?ixlib=rb-1.2.1&auto=format&fit=crop&w=500&q=60",
		},
		{
			Name:        "IKEA 台灯",
			Description: "毕业带不走，低价出。",
			Price:       25.00,
			CategoryId:  category.Id,
			SellerId:    users[0].ID,
			Status:      "available",
			Condition:   "fair",
			ImageUrl:    "https://images.unsplash.com/photo-1534234828569-1f27c78ee755?q=80&w=500&auto=format&fit=crop",
		},
		{
			Name:        "考研数学复习全书",
			Description: "买了没怎么看，99新，附赠笔记。",
			Price:       15.00,
			CategoryId:  categoryBook.Id,
			SellerId:    users[1].ID, // 用户2发布的商品
			Status:      "sold",      // 已售出
			Condition:   "new",
			ImageUrl:    "https://images.unsplash.com/photo-1544716278-ca5e3f4abd8c?w=500&auto=format&fit=crop",
		},
	}
	if err := db.Create(&products).Error; err != nil {
		log.Println(err)
	}
	fmt.Printf("✅ 商品创建完成: %d 个\n", len(products))

	// 7. 创建订单与评价 (Mock 核心逻辑)
	seedOrdersAndReviews(db, users, products)

	fmt.Println("🎉 数据播种全部完成！")
}

func seedOrdersAndReviews(db *gorm.DB, users []model.User, products []model.Product) {
	// 场景 1: 用户2 购买了 用户1 的 Sony 耳机 (products[1])
	order1 := model.Order{
		ProductId:   products[1].Id,
		BuyerId:     users[1].ID,
		SellerId:    users[0].ID,
		Status:      "completed",
		Amount:      products[1].Price,
		CompletedAt: time.Now().Add(-24 * time.Hour), // 昨天完成的
	}
	if err := db.Create(&order1).Error; err != nil {
		panic(err)
	}

	// 评价 1: 用户2 -> 用户1
	review1 := model.Review{
		OrderId:      order1.Id,
		ReviewerId:   users[1].ID, // 评价人：买家
		TargetUserId: users[0].ID, // 被评价人：卖家
		Rating:       5,
		Comment:      "耳机音质真的绝了，卖家包装也很用心，还送了贴纸，好评！",
		CreatedAt:    time.Now(),
	}
	if err := db.Create(&review1).Error; err != nil {
		panic(err)
	}

	// 场景 2: 用户1 购买了 用户2 的 考研书 (products[3])
	order2 := model.Order{
		ProductId:   products[3].Id,
		BuyerId:     users[0].ID,
		SellerId:    users[1].ID,
		Status:      "completed",
		Amount:      products[3].Price,
		CompletedAt: time.Now().Add(-48 * time.Hour),
	}
	if err := db.Create(&order2).Error; err != nil {
		panic(err)
	}

	// 评价 2: 用户1 -> 用户2
	review2 := model.Review{
		OrderId:      order2.Id,
		ReviewerId:   users[0].ID,
		TargetUserId: users[1].ID,
		Rating:       4,
		Comment:      "书是正版的，笔记也很详细，就是快递稍微慢了点。",
		CreatedAt:    time.Now(),
	}
	if err := db.Create(&review2).Error; err != nil {
		panic(err)
	}

	fmt.Println("✅ 订单与评价 Mock 完成")
}

func cleanData(db *gorm.DB) {
	// 注意删除顺序，避免外键约束报错
	db.Exec("DELETE FROM reviews")
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM product_sold_logs")
	db.Exec("DELETE FROM product_drop_logs")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM categories")
	db.Exec("DELETE FROM users")

	// 重置自增 ID (MySQL)
	tables := []string{"reviews", "orders", "products", "categories", "users", "product_sold_logs"}
	for _, t := range tables {
		db.Exec(fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = 1", t))
	}
}

func dropData(db *gorm.DB) {
	db.Exec("DROP TABLE products")
	db.Exec("DROP TABLE categories")
	db.Exec("DROP TABLE users")
}
