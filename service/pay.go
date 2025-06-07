package service

import (
	"errors"
	"parking-system-go/global"
	"parking-system-go/model/database"
)

type PayService struct{}

// 统一接口实现创建，支付，查询，
func (s *PayService) Pay(req *database.Order) error {

	if err := s.CreateMockOrder(req); err != nil {
		return err
	}
	s.MockPayment(req.OrderID)

	return nil
}

func (s *PayService) CreateMockOrder(req *database.Order) error {
	var order database.Order
	err := global.DB.Where("order_id = ?", req.OrderID).First(&order).Error
	if err == nil {
		return errors.New("订单已创建")
	}

	if err := global.DB.Create(req).Error; err != nil {
		return errors.New("创建订单失败")
	}
	return nil
}

func (s *PayService) MockPayment(orderID string) error {
	var order database.Order
	err := global.DB.Where("order_id = ?", orderID).First(&order).Error
	if err != nil {
		return errors.New("订单不存在")
	}

	order.Status = 1
	if err := global.DB.Save(&order).Error; err != nil {
		return errors.New("设置支付失败")
	}
	return nil
}

func (s *PayService) GetOrderStatus(orderID string) (int8, error) {
	var order database.Order
	err := global.DB.Where("order_id = ?", orderID).First(&order).Error
	if err != nil {
		return 0, errors.New("订单不存在")
	}
	return order.Status, nil
}
