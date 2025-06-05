package api

import "parking-system-go/service"

type ApiGroup struct {
	ParkingApi
	BarrierLogApi
	UserApi
}

var ApiGroupApp = new(ApiGroup)
var parkingService = service.ServiceGroupApp.ParkingService
var userService = service.ServiceGroupApp.UserService
var barrierLogService = service.ServiceGroupApp.BarrierLogService
var parkingrecordService = service.ServiceGroupApp.ParkingRecordService
var parkinglotService = service.ServiceGroupApp.ParkingLotService
