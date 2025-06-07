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

	// 参数检查
	if req.ID == 0 && req.ParkingRecordID == 0 && req.OrderID == "" {
		return order, errors.New("必须提供订单编号或停车记录 ID")
	}

	query := global.DB.Model(&database.Order{})

	// 依优先级判断查询条件
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	} else if req.OrderID != "" {
		query = query.Where("order_id = ?", req.OrderID)
	} else if req.ParkingRecordID != 0 {
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
