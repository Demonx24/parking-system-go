package service

type ServiceGroup struct {
	ParkingService
	UserService
}

var ServiceGroupApp = new(ServiceGroup)
