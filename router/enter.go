package router

type RouterGroup struct {
	ParkingRouter
	BarrierLogRouter
	UserRouter
	PayRouter
}

var RouterGroupApp = new(RouterGroup)
