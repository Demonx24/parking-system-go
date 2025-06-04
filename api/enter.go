package api

import "parking-system-go/service"

type ApiGroup struct {
	ParkingApi
}

var ApiGroupApp = new(ApiGroup)
var parkingService = service.ServiceGroupApp.ParkingService
var userService = service.ServiceGroupApp.UserService
