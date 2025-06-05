package router

type RouterGroup struct {
	ParkingRouter
	BarrierLogRouter
	UserRouter
}

var RouterGroupApp = new(RouterGroup)
