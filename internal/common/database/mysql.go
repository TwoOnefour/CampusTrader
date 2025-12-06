package database

import (
	"CampusTrader/internal/model"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局变量，其他包直接用 database.DB 调用
var DB *gorm.DB

func InitMySQL() {
	// 1. 配置 DSN (Data Source Name)
	// 格式: 用户名:密码@tcp(IP:端口)/数据库名?配置项
	// parseTime=True: 自动把数据库时间转为 Go 的 time.Time
	// loc=Local: 使用本地时区
	dsn := os.Getenv("MYSQL_DSN")

	// 2. 配置 GORM 日志 (这对于毕设调试非常重要！)
	// 这样控制台会打印出每一条执行的 SQL 语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 级别：Info 会打印所有 SQL
			IgnoreRecordNotFoundError: true,        // 忽略 ErrRecordNotFound 错误
			Colorful:                  true,        // 彩色打印
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

	fmt.Println("MySQL 连接成功！")
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		return
	}
}
