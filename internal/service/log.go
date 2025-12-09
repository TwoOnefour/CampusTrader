package service

import (
	"CampusTrader/internal/model"
	"gorm.io/gorm"
)

type LogService struct {
	db *gorm.DB
}

func NewLogService(db *gorm.DB) *LogService {
	return &LogService{
		db: db,
	}
}

func (s *LogService) OnOrder(order *model.Order) error {
	db := s.db.Model(&model.ProductSoldLog{})

	return db.Create(&model.ProductSoldLog{
		ProductId: order.ProductId,
		BuyerId:   order.BuyerId,
		SellerId:  order.SellerId,
		Price:     order.Amount,
	}).Error
}

// 不要用gin 的 ctx
func (s *LogService) OnProductDrop(productId, operatorId uint64, reason string) error {
	db := s.db.Model(&model.ProductDropLogs{})
	return db.Create(&model.ProductDropLogs{
		OperatorId: operatorId,
		ProductId:  productId,
		Reason:     reason,
	}).Error
}
