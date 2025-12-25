package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	Id          uint64         `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;" json:"id"`
	Name        string         `gorm:"column:name;type:VARCHAR(100);comment:商品名称;not null;" json:"name"`                  // 商品名称
	Description string         `gorm:"column:description;type:TEXT;comment:商品描述;not null;" json:"description"`            // 商品描述
	Price       float64        `gorm:"column:price;type:DECIMAL(10, 2);comment:价格;not null;" json:"price"`                // 价格
	CategoryId  uint64         `gorm:"column:category_id;type:BIGINT UNSIGNED;comment:分类ID;not null;" json:"category_id"` // 分类ID
	SellerId    uint64         `gorm:"column:seller_id;type:BIGINT UNSIGNED;comment:卖家ID;not null;" json:"seller_id"`     // 卖家ID
	Category    Category       `gorm:"foreignKey:CategoryId" json:"category,omitempty"`
	Seller      User           `gorm:"foreignKey:SellerId;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"seller"`
	Status      string         `gorm:"column:status;type:ENUM('available', 'sold', 'removed');comment:状态;default:available;" json:"status"`           // 状态
	Condition   string         `gorm:"column:condition;type:ENUM('new', 'like_new', 'good', 'fair', 'poor');comment:新旧程度;not null;" json:"condition"` // 新旧程度
	ImageUrl    string         `gorm:"column:image_url;type:VARCHAR(255);comment:主图URL;" json:"image_url"`                                            // 主图URL
	CreatedAt   time.Time      `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type ProductWithUserRating struct {
	Product
	RatingStat RatingStat `json:"user_rating_stat"`
}

type ProductSoldLog struct {
	Id        uint64    `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;" json:"id"`
	ProductId uint64    `gorm:"column:product_id;type:BIGINT UNSIGNED;not null;" json:"product_id"`
	BuyerId   uint64    `gorm:"column:buyer_id;type:BIGINT UNSIGNED;not null;" json:"buyer_id"`
	SellerId  uint64    `gorm:"column:seller_id;type:BIGINT UNSIGNED;not null;" json:"seller_id"`
	Price     float64   `gorm:"column:price;type:DECIMAL(10, 2);comment:成交价;not null;" json:"price"` // 成交价
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"created_at"`
}

type ProductDropLogs struct {
	Id         uint64    `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;" json:"id"`
	ProductId  uint64    `gorm:"column:product_id;type:BIGINT UNSIGNED;not null;" json:"product_id"`
	OperatorId uint64    `gorm:"column:operator_id;type:BIGINT UNSIGNED;comment:操作人ID;not null;" json:"operator_id"` // 操作人ID
	Reason     string    `gorm:"column:reason;type:VARCHAR(255);comment:下架原因;" json:"reason"`                        // 下架原因
	CreatedAt  time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"created_at"`
}

type PageParam struct {
	LastId   uint64 `form:"last_id,default=0"`
	PageSize uint64 `form:"page_size,default=50"`
}

type PageData[T any] struct {
	List    []T    `json:"list"`
	HasMore bool   `json:"has_more"`
	LastId  uint64 `json:"last_id"`
}
