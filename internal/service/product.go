package service

import (
	"CampusTrader/internal/model"
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type ProductService struct {
	db          *gorm.DB
	logService  *LogService
	statService *StatisticsService
}

func NewProductService(db *gorm.DB, logService *LogService, statService *StatisticsService) *ProductService {
	return &ProductService{
		db:          db,
		logService:  logService,
		statService: statService,
	}
}

func (s *ProductService) CreateProduct(ctx context.Context, req *model.Product) error {
	db := s.db.WithContext(ctx)
	var count int64
	db.Model(&model.Product{}).
		Where("user_id = ? AND title = ? AND price = ? AND created_at > ?",
			req.SellerId, req.Name, req.Price, time.Now().Add(-5*time.Second)).
		Count(&count)

	if count > 0 {
		return errors.New("请勿重复提交，休息一下吧")
	}
	if err := db.Create(req).Error; err != nil {
		return err
	}
	return nil
}

func (s *ProductService) ListProducts(ctx context.Context, pageSize, lastID uint64) ([]model.ProductWithUserRating, error) {
	db := s.db.WithContext(ctx).Model(&model.Product{}).Where("status = ?", "available")
	if lastID > 0 {
		db = db.Where("id < ?", lastID)
	}
	var products []model.Product
	err := db.Order("id DESC").
		Limit(int(pageSize)).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	UserIDs := make([]uint64, 0)
	UserIDSet := make(map[uint64]struct{}) // 去重
	for _, p := range products {
		if _, ok := UserIDSet[p.SellerId]; ok {
			continue
		}
		UserIDs = append(UserIDs, p.SellerId)
		UserIDSet[p.SellerId] = struct{}{}
	}
	//sql := `
	//	create or replace view v_product as
	//		select
	//			products.*, users.nickname,
	//			users.phone,
	//			IFNULL(t.avg_rating, 0),
	//			IFNULL(t.review_count, 0)
	//		from products
	//		join users on products.seller_id = users.id
	//		LEFT JOIN (
	//			SELECT
	//				target_user_id,
	//				AVG(rating) AS avg_rating,
	//				COUNT(*) AS review_count
	//			FROM reviews
	//			GROUP BY target_user_id
	//		) t ON t.target_user_id = users.id;
	//	select * from v_product;
	//`
	// 这里可以用这个view，但我直接分查了
	ratings, err := s.statService.BatchGetUserRating(ctx, UserIDs)
	if err != nil {
		return nil, err
	}
	productsWithRatings := make([]model.ProductWithUserRating, len(products))
	for i, p := range products {
		productsWithRatings[i] = model.ProductWithUserRating{
			Product: p,
		}
		if stat, ok := ratings[p.SellerId]; ok {
			productsWithRatings[i].RatingStat = stat
		}
	}
	return productsWithRatings, err
}

func (s *ProductService) ListMyProducts(ctx context.Context, sellerID uint64) ([]model.ProductWithUserRating, error) {
	var products []model.Product
	err := s.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("seller_id = ?", sellerID).
		Order("created_at DESC"). // 按时间倒序
		Find(&products).Error
	rating, err := s.statService.GetUserRating(ctx, sellerID)
	if err != nil {
		return nil, err
	}
	productsWithRatings := make([]model.ProductWithUserRating, len(products))
	for i, p := range products {
		productsWithRatings[i] = model.ProductWithUserRating{
			Product:    p,
			RatingStat: *rating,
		}
	}
	return productsWithRatings, err
}

func (s *ProductService) SearchProduct(ctx context.Context, keyword string, count uint64) ([]model.Product, error) {
	db := s.db.WithContext(ctx).Model(&model.Product{})
	var products []model.Product
	err := db.Where("name like ?", "%"+keyword+"%").Order("id desc").Limit(int(count)).Find(&products).Error
	return products, err
}

func (s *ProductService) SearchSuggestion(ctx context.Context, keyword string) ([]string, error) {
	db := s.db.WithContext(ctx).Model(&model.Product{})
	var products []model.Product
	err := db.Where("name like ? and status = 'available'", "%"+keyword+"%").Order("id desc").Limit(5).Find(&products).Error
	if err != nil {
		return nil, err
	}
	productPrefix := make([]string, len(products))
	for i := 0; i < len(products); i++ {
		productPrefix[i] = products[i].Name
	}
	return productPrefix, err
}

func (s *ProductService) DropProduct(ctx context.Context, productId, operatorId uint64, reason string) error {
	db := s.db.WithContext(ctx)
	err := db.Model(&model.Product{}).Delete(&model.Product{
		Id: productId,
	}, productId).Error
	if err != nil {
		return err
	}
	notifyLogging := func() {
		err := s.logService.OnProductDrop(productId, operatorId, reason)
		if err != nil {
			log.Println(err.Error())
		}
	}
	go notifyLogging()
	return nil
}
