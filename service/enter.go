package service

type ServiceGroup struct {
	ParkingService
	UserService
	BarrierLogService
	ParkingLotService
	ParkingRecordService
	OrderService
	PayService
	PayWeChatService
}

var ServiceGroupApp = new(ServiceGroup)
