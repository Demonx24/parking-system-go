package api

import (
	"github.com/gin-gonic/gin"
)

type UserApi struct{}

func (userApi *UserApi) GetUser(c *gin.Context) {
	//var req database.User
	//if err := c.ShouldBindJSON(&req); err != nil {
	//	global.Log.Error("查询用户信息参数错误", zap.Error(err))
	//}
	//if user, err := userService.GetUser(req); err != nil {
	//	global.Log.Error("查询用户信息数据失败", zap.Error(err))
	//}

}
