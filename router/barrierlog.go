package router

import (
	"github.com/gin-gonic/gin"
	"parking-system-go/api"
)

type BarrierLogRouter struct{}

func (b *BarrierLogRouter) InitBarrierLogRouter(Router *gin.RouterGroup) {
	barrierLogRouter := Router.Group("barrierlog")
	barrierLogApi := api.ApiGroupApp.BarrierLogApi
	{
		barrierLogRouter.POST("create", barrierLogApi.BarrierLogCreate)
	}
}
