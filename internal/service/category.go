package service

import (
	"CampusTrader/internal/model"
	"context"
	"gorm.io/gorm"
)

type CategoryService struct {
	db *gorm.DB
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{
		db: db,
	}
}

// 懒得写流式了
func (s *CategoryService) ListCategory(ctx context.Context) ([]model.Category, error) {
	db := s.db.WithContext(ctx)
	var categories []model.Category
	db.Model(&model.Category{}).Find(&categories)
	return categories, nil
}

func (s *CategoryService) CreateCategory(ctx context.Context, category model.Category) error {
	db := s.db.WithContext(ctx)
	return db.Create(&category).Error
}

// 查询递归所有父类
func (s *CategoryService) ListRelatedCategory(ctx context.Context, id uint) ([]model.Category, error) {
	db := s.db.WithContext(ctx).Model(&model.Category{})
	var categories []model.Category
	sql := `
        WITH RECURSIVE CategoryPath AS (
            SELECT id, name, parent_id, created_at, updated_at
            FROM categories
            WHERE id = ?
            
            UNION ALL
            
            SELECT c.id, c.name, c.parent_id, c.created_at, c.updated_at
            FROM categories c
            INNER JOIN CategoryPath cp ON cp.parent_id = c.id
        )
        SELECT * FROM CategoryPath;
    `
	db.Raw(sql, id).Scan(&categories)
	return categories, nil
}
