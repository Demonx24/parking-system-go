package router

import (
	"github.com/gin-gonic/gin"
	"parking-system-go/api"
)

type UserRouter struct{}

func (b *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	usereRouter := Router.Group("user")
	UserApi := api.ApiGroupApp.UserApi
	{
		usereRouter.POST("create", UserApi.CreateUser)
		usereRouter.GET("get", UserApi.GetUser)
	}
}
