package database

import (
	"CampusTrader/internal/model"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitMySQL() {
	// 1. 配置 DSN (Data Source Name)
	// 格式: 用户名:密码@tcp(IP:端口)/数据库名?配置项
	// parseTime=True: 自动把数据库时间转为 Go 的 time.Time
	// loc=Local: 使用本地时区
	dsn := os.Getenv("DATABASE_DSN")

	// 2. 配置 GORM 日志 (这对于毕设调试非常重要！)
	// 这样控制台会打印出每一条执行的 SQL 语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  logger.Error, // 级别：Info 会打印所有 SQL
			IgnoreRecordNotFoundError: true,         // 忽略 ErrRecordNotFound 错误
			Colorful:                  true,         // 彩色打印
		},
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	// 3. 配置连接池 (Connection Pool)
	// 获取底层的 sql.DB 对象
	sqlDB, err := DB.DB()
	if err != nil {
		panic("获取底层 sql.DB 失败")
	}

	sqlDB.SetMaxIdleConns(10)

	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	allModels := []interface{}{
		&model.User{},
		&model.Product{},
		&model.Category{},
		&model.Order{},
		&model.ProductSoldLog{},
		&model.ProductDropLogs{},
		&model.Review{},
	}
	if err := DB.AutoMigrate(allModels...); err != nil {
		return
	}
	initAdvancedObjects(DB)
}

func initAdvancedObjects(db *gorm.DB) {
	viewProductDetails := `
	CREATE OR REPLACE VIEW v_product_details AS
	SELECT 
		p.*, 
		c.name AS category_name, 
		u.nickname AS seller_nickname
	FROM products p
	JOIN categories c ON p.category_id = c.id
	JOIN users u ON p.seller_id = u.id;`
	viewUserTradeSummary := `
	CREATE OR REPLACE VIEW v_user_trade_summary AS
	SELECT 
		u.id, u.username, u.nickname,
		(SELECT COUNT(*) FROM products WHERE seller_id = u.id) AS total_listed,
		(SELECT SUM(amount) FROM orders WHERE seller_id = u.id AND status = 'completed') AS total_sales_amount
	FROM users u;`
	db.Exec(`DROP PROCEDURE IF EXISTS sp_search_and_count_by_category;`)
	db.Exec(`DROP PROCEDURE IF EXISTS sp_complete_order;`)
	procCompleteOrder := `
	CREATE PROCEDURE sp_complete_order(IN p_order_id BIGINT UNSIGNED)
	BEGIN
		UPDATE orders SET status = 'completed', updated_at = NOW() WHERE id = p_order_id;
		UPDATE products SET status = 'sold' WHERE id = (SELECT product_id FROM orders WHERE id = p_order_id);
	END;`

	procSearchCount := `
	CREATE PROCEDURE sp_search_and_count_by_category(
		IN p_category_id BIGINT UNSIGNED
	)
	BEGIN
		SELECT * FROM products
		where 
			category_id=p_category_id 
			AND status='available';
	END;`

	objects := []string{viewProductDetails, viewUserTradeSummary, procCompleteOrder, procSearchCount}
	for _, sql := range objects {
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("创建数据库对象失败: %v", err)
		}
	}
}
