package api

import "parking-system-go/service"

type ApiGroup struct {
	ParkingApi
	BarrierLogApi
	UserApi
	PayApi
	PayWeChatApi
}

var ApiGroupApp = new(ApiGroup)
var parkingService = service.ServiceGroupApp.ParkingService
var userService = service.ServiceGroupApp.UserService
var barrierLogService = service.ServiceGroupApp.BarrierLogService
var parkingrecordService = service.ServiceGroupApp.ParkingRecordService
var parkinglotService = service.ServiceGroupApp.ParkingLotService
var payService = service.ServiceGroupApp.PayService
var orderService = service.ServiceGroupApp.OrderService
var payWeChatService = service.ServiceGroupApp.PayWeChatService
