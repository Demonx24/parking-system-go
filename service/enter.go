package service

type ServiceGroup struct {
	ParkingService
	UserService
	BarrierLogService
	ParkingLotService
	ParkingRecordService
	OrderService
}

var ServiceGroupApp = new(ServiceGroup)
