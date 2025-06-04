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
	}
}
