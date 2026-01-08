package database

import "gorm.io/gorm"

// DB 全局变量，其他包直接用 database.DB 调用
var DB *gorm.DB
