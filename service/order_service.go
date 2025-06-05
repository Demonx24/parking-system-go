package service

import (
	"errors"
	"parking-system-go/global"
	"parking-system-go/model/database"
)

type OrderService struct{}

func (s *OrderService) Create(order *database.Order) error {
	return global.DB.Create(order).Error
}

func (s *OrderService) GetOrder(req database.Order) (database.Order, error) {
	var order database.Order
	if req.ID == 0 && req.ParkingRecordID == 0 {
		return order, errors.New("必须提供 ID 或停车记录 ID")
	}

	query := global.DB.Model(&database.Order{})
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	} else {
		query = query.Where("parking_record_id = ?", req.ParkingRecordID)
	}

	err := query.First(&order).Error
	if err != nil {
		return order, err
	}
	return order, nil
}

func (s *OrderService) Update(order *database.Order) error {
	return global.DB.Save(order).Error
}

func (s *OrderService) Delete(where *database.Order) error {
	return global.DB.Where(where).Delete(&database.Order{}).Error
}
