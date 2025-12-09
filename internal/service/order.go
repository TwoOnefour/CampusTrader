package service

import (
	"CampusTrader/internal/model"
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
)

type OrderService struct {
	db         *gorm.DB
	logService *LogService
}

func NewOrderService(db *gorm.DB, logService *LogService) *OrderService {
	return &OrderService{
		db:         db,
		logService: logService,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, itemID, buyerID uint64) error {
	var item model.Product
	db := s.db.WithContext(ctx)
	// TODO: redis lua
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
			SellerId:  item.SellerId,
		}).Error
	})

	if err != nil {
		return err
	}

	go func() {
		err := s.logService.OnOrder(&model.Order{
			ProductId: item.Id,
			BuyerId:   buyerID,
			Amount:    item.Price,
			SellerId:  item.SellerId,
		})
		if err != nil {
			log.Println(err.Error())
		}
	}()

	return nil
}

func (s *OrderService) ListOrder(ctx context.Context, userId uint) ([]model.Order, error) {
	db := s.db.WithContext(ctx)
	var order []model.Order
	db.Model(&model.Order{}).Where("buyer_id=?", userId).Find(&order)
	return order, nil
}
