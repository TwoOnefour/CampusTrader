package page

import (
	"CampusTrader/internal/model"

	"gorm.io/gorm"
)

type Result[T any] struct {
	List    []T    `json:"list"`
	HasMore bool   `json:"has_more"`
	LastId  uint32 `json:"last_id"` // 假设你的 ID 统一是 uint32
}

func PaginateExec[T any](db *gorm.DB, p model.PageParam) (*Result[T], error) {
	var list []T
	pageSize := int(p.PageSize)

	// 1. 执行查询：多查一条 (pageSize + 1) 用于判断 HasMore
	// 这里的 db 应该是已经包含了各种 Where 条件的基础查询
	err := db.Order("id DESC").Limit(pageSize + 1).Find(&list).Error
	if err != nil {
		return nil, err
	}

	res := &Result[T]{
		List:    list,
		HasMore: false,
	}

	// 2. 逻辑判断：是否有下一页
	if len(list) > pageSize {
		res.HasMore = true
		res.List = list[:pageSize] // 截断到前端申请的长度
	}

	return res, nil
}
