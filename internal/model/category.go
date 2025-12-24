package model

import "time"

type Category struct {
	Id        uint64    `gorm:"column:id;type:BIGINT UNSIGNED;primaryKey;" json:"id"`
	Name      string    `gorm:"column:name;type:VARCHAR(50);comment:分类名称;not null;" json:"name"`                    // 分类名称
	ParentId  uint64    `gorm:"column:parent_id;type:BIGINT UNSIGNED;comment:父分类ID;default:NULL;" json:"parent_id"` // 父分类ID
	Parent    *Category `gorm:"foreignKey:ParentId" json:"parent,omitempty"`
	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"updated_at"`
}
