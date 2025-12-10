package service

import (
	"CampusTrader/internal/model"
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type ProductService struct {
	db         *gorm.DB
	logService *LogService
}

func NewProductService(db *gorm.DB, logService *LogService) *ProductService {
	return &ProductService{
		db:         db,
		logService: logService,
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

func (s *ProductService) ListProducts(ctx context.Context, pageSize, lastID uint64) ([]model.Product, error) {
	db := s.db.WithContext(ctx).Model(&model.Product{}).Where("status = ?", "available")
	if lastID > 0 {
		db = db.Where("id < ?", lastID)
	}
	var products []model.Product
	err := db.Order("id DESC").
		Limit(int(pageSize)).
		Find(&products).Error
	return products, err
}

func (s *ProductService) ListMyProducts(ctx context.Context, sellerID uint64) ([]model.Product, error) {
	var products []model.Product
	err := s.db.WithContext(ctx).
		Model(&model.Product{}).
		Where("seller_id = ?", sellerID).
		Order("created_at DESC"). // 按时间倒序
		Find(&products).Error
	return products, err
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
