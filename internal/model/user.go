package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"user_id"`
	Username  string    `gorm:"column:username;type:VARCHAR(50);not null;unique;index;" json:"username"`
	Password  string    `gorm:"column:password;type:VARCHAR(255);not null;index;" json:"password"`
	Nickname  string    `gorm:"column:nickname;type:VARCHAR(50);not null;" json:"nickname"`
	Email     string    `gorm:"column:email;type:VARCHAR(100);not null;unique;index;" json:"email"`
	Phone     string    `gorm:"column:phone;type:VARCHAR(20);" json:"phone"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"created_at"`
}
