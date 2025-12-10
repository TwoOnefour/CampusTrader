package service

import (
	"CampusTrader/internal/model"
	"context"
	"gorm.io/gorm"
)

type StatisticsService struct {
	db *gorm.DB
}

func NewStatisticsService(db *gorm.DB) *StatisticsService {
	return &StatisticsService{db: db}
}

// 复杂查询最热门种类TOP3, 思路为先连两个表, 根据目录id分组计数降序排序取top3，然后根据id连回categories表查名字等信息
func (s *StatisticsService) GetHotCategories(ctx context.Context) ([]model.Category, error) {
	db := s.db.WithContext(ctx)
	var categories []model.Category
	sql := `
		select
			categories.name as name,
			t.category_id as id
		from categories
			inner join
				(
					SELECT
						products.category_id,
						count(*) as cnt
					FROM orders
						INNER JOIN products
							ON products.id = orders.product_id
					where orders.status in ('paid', 'completed')
					group by products.category_id
					order by cnt desc
					limit 3
				) as t
			on t.category_id = categories.id
		order by t.cnt
	`
	if err := db.Raw(sql).Scan(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (s *StatisticsService) GetUserRating(ctx context.Context, userId uint64) (float64, error) {
	type Result struct {
		Rating float64
	}
	var res Result
	// 很坑，我本来写的sum(if(rating > 3, 1, 0) / count(*)，但该用户没有评价记录的时候会返回null，导致gorm报错，只能换成下面这样
	sql := `
       SELECT 
          IFNULL(AVG(rating > 3), 0) AS rating 
       FROM reviews 
       WHERE target_user_id = ?
    `

	if err := s.db.WithContext(ctx).Raw(sql, userId).Scan(&res).Error; err != nil {
		return 0, err
	}
	
	return res.Rating, nil
}

func (s *StatisticsService) GetUserCompletedOrderRecord(ctx context.Context, userId uint64) ([]model.Order, error) {
	var orders []model.Order
	db := s.db.WithContext(ctx).Model(&model.Order{})
	if err := db.Where("(seller_id = ? or buyer_id = ?)", userId, userId).Where("status = 'completed'").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
