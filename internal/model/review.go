package model

import "time"

type Review struct {
	Id           uint64    `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;" json:"id"`
	OrderId      uint64    `gorm:"column:order_id;type:BIGINT UNSIGNED;not null;" json:"order_id"`
	ReviewerId   uint64    `gorm:"column:reviewer_id;type:BIGINT UNSIGNED;comment:评价人;not null;" json:"reviewer_id"`        // 评价人
	TargetUserId uint64    `gorm:"column:target_user_id;type:BIGINT UNSIGNED;comment:被评价人;not null;" json:"target_user_id"` // 被评价人
	Rating       uint8     `gorm:"column:rating;type:TINYINT UNSIGNED;comment:评分1-5;not null;" json:"rating"`               // 评分1-5
	Comment      string    `gorm:"column:comment;type:TEXT;comment:评价内容;" json:"comment"`                                   // 评价内容
	CreatedAt    time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"created_at"`
}

type RatingStat struct {
	TargetUserID uint64  `gorm:"column:target_user_id" json:"target_user_id"`
	AvgRating    float64 `gorm:"column:avg_rating" json:"avg_rating"`
	ReviewCount  uint64  `gorm:"column:review_count" json:"review_count"`
}
