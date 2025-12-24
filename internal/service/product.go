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

func paginate(p model.PageParam) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.LastId > 0 {
			db = db.Where("id < ?", p.LastId)
		}
		return db.Order("id DESC")
	}
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

func (s *ProductService) ListProducts(ctx context.Context, pageParam model.PageParam) (*model.PageData[model.ProductWithUserRating], error) {
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
	db := s.db.WithContext(ctx).Model(&model.Product{}).Where("status = ?", "available")
	db = paginate(pageParam)(db)
	var products []model.Product
	err := db.Limit(int(pageParam.PageSize + 1)).Find(&products).Error
	if err != nil {
		return nil, err
	}
	type resType = model.PageData[model.ProductWithUserRating]
	var res resType
	productWithRating, err := getUserRatingsByProducts(ctx, s.statService, products)
	if err != nil {
		return nil, err
	}
	res = resType{
		List:    productWithRating,
		HasMore: true,
	}
	return &res, nil
}

func (s *ProductService) ListMyProducts(ctx context.Context, sellerID uint64, pageParam model.PageParam) (*model.PageData[model.ProductWithUserRating], error) {
	var products []model.Product
	db := s.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("seller_id = ?", sellerID).
		Order("created_at DESC")
	db = paginate(pageParam)(db)
	if err := db.Find(&products).Error; err != nil {
		return nil, err
	}
	type resType = model.PageData[model.ProductWithUserRating]
	var res resType
	productWithRating, err := getUserRatingsByProducts(ctx, s.statService, products)
	if err != nil {
		return nil, err
	}
	res = resType{
		List:    productWithRating,
		HasMore: true,
	}
	return &res, nil
}

func (s *ProductService) ListProductsByProc(ctx context.Context, categoryID uint64, pageParam model.PageParam) (*model.PageData[model.ProductWithUserRating], error) {
	var products []model.Product
	db := s.db.WithContext(ctx)
	db = paginate(pageParam)(db)
	err := db.Raw("CALL sp_search_and_count_by_category(?)", categoryID).Scan(&products).Error
	if err != nil {
		return nil, err
	}
	type resType = model.PageData[model.ProductWithUserRating]
	var res resType
	productWithRating, err := getUserRatingsByProducts(ctx, s.statService, products)
	if err != nil {
		return nil, err
	}
	res = resType{
		List:    productWithRating,
		HasMore: true,
	}
	return &res, nil
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

func getUserRatingsByProducts(ctx context.Context, s *StatisticsService, products []model.Product) ([]model.ProductWithUserRating, error) {
	UserIDs := make([]uint64, 0)
	UserIDSet := make(map[uint64]struct{}) // 去重
	for _, p := range products {
		if _, ok := UserIDSet[p.SellerId]; ok {
			continue
		}
		UserIDs = append(UserIDs, p.SellerId)
		UserIDSet[p.SellerId] = struct{}{}
	}
	ratings, err := s.BatchGetUserRating(ctx, UserIDs)
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
	return productsWithRatings, nil
}
