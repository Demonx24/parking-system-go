package router

import (
	"github.com/gin-gonic/gin"
	"parking-system-go/api"
)

type ParkingRouter struct{}

func (p *ParkingRouter) InitParkingRouter(Router *gin.RouterGroup) {
	parkingRouter := Router.Group("parking")
	parkingApi := api.ApiGroupApp.ParkingApi
	{
		parkingRouter.GET("status", parkingApi.GetParkingStatus)
		parkingRouter.POST("entry", parkingApi.Entry)
		parkingRouter.POST("exit", parkingApi.Exit)
		//parkingRouter.POST("create_payment", parkingApi.CreatePayment)
		//parkingRouter.POST("payment_notify", parkingApi.PaymentNotify)
	}
}
