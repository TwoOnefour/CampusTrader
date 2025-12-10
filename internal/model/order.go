package model

import "time"

type Order struct {
	Id          uint64    `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;" json:"id"`
	ProductId   uint64    `gorm:"column:product_id;type:BIGINT UNSIGNED;comment:商品ID;not null;" json:"product_id"` // 商品ID
	BuyerId     uint64    `gorm:"column:buyer_id;type:BIGINT UNSIGNED;comment:买家ID;not null;" json:"buyer_id"`     // 买家ID
	SellerId    uint64    `gorm:"column:seller_id;type:BIGINT UNSIGNED;comment:卖家ID;not null;" json:"seller_id"`
	Status      string    `gorm:"column:status;type:ENUM('pending', 'paid', 'completed', 'cancelled');comment:订单状态;default:pending;" json:"status"` // 订单状态
	Amount      float64   `gorm:"column:amount;type:DECIMAL(10, 2);comment:实际成交金额;not null;" json:"amount"`                                         // 实际成交金额
	CreatedAt   time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	CompletedAt time.Time `gorm:"column:completed_at;type:TIMESTAMP;comment:完成时间;default:CURRENT_TIMESTAMP;" json:"completed_at"` // 完成时间
}
