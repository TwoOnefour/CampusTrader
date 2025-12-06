package service

import (
	"CampusTrader/internal/model"
	"context"
	"errors"
	"gorm.io/gorm"
)

type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{
		db: db,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, itemID, buyerID uint64) error {
	var item model.Product
	db := s.db.WithContext(ctx)
	db.First(&item, itemID)
	if item.Status != "available" {
		return errors.New("已售出")
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		// 乐观锁
		res := tx.Model(&model.Product{}).Where("id=? AND status='available'", itemID).Update("status", "sold")
		if res.RowsAffected == 0 {
			return errors.New("已售出")
		}
		return tx.Create(&model.Order{
			ProductId: item.Id,
			BuyerId:   buyerID,
			Amount:    item.Price,
		}).Error
	})

	if err != nil {
		return err
	}

	// go s.SendNotification(...)

	return nil
}
