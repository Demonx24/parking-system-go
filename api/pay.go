package api

import (
	"github.com/gin-gonic/gin"
	"parking-system-go/global"
	"parking-system-go/model/database"
	"parking-system-go/model/response"
)

// 接口以及写成方法，在service层
type PayApi struct{}

func (p *PayApi) CreateMockOrder(c *gin.Context) {
	var req database.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}
	order, _ := orderService.GetOrder(req)
	//if err != nil {
	//	response.FailWithMessage("查询订单失败", c)
	//	return
	//}
	if order.OrderID != "" {
		response.FailWithMessage("订单以创建", c)
		return
	}
	if err := orderService.Create(&req); err != nil {
		response.FailWithMessage("创建订单失败", c)
		return
	}
	response.OkWithData(req, c)
}

func (p *PayApi) MockPayment(c *gin.Context) {
	var req database.Order
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误", c)
		return
	}

	req, err := orderService.GetOrder(req)
	if err != nil {
		response.FailWithMessage("没有该订单", c)
		return
	}
	req.Status = 1
	if err := orderService.Update(&req); err != nil {
		response.FailWithMessage("设置支付失败", c)
		return
	}
	response.OkWithMessage("模拟支付成功", c)
}
func (p *PayApi) GetOrderStatus(c *gin.Context) {
	orderID := c.Query("order_id")
	var order database.Order
	if err := global.DB.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		response.FailWithMessage("订单不存在", c)
		return
	}
	response.OkWithData(gin.H{"status": order.Status}, c)
}
