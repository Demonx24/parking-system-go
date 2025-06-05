package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"parking-system-go/global"
	"parking-system-go/model/database"
	"parking-system-go/model/response"
)

type UserApi struct{}

func (userApi *UserApi) CreateUser(c *gin.Context) {
	var req database.User
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("查询用户信息参数错误", zap.Error(err))
	}
	if err := userService.Update(&req); err != nil {
		global.Log.Error("<UNK>", zap.Error(err))
		response.FailWithMessage("<UNK>", c)
	}
	response.OkWithData("<UNK>", c)
}
func (userApi *UserApi) GetUser(c *gin.Context) {
	var req database.User
	if err := c.ShouldBindQuery(&req); err != nil {
		global.Log.Error("查询用户信息参数错误", zap.Error(err))
		response.FailWithMessage("<UNK>", c)
		return
	}
	user, err := userService.GetUser(req)
	if err != nil {
		global.Log.Error("<UNK>", zap.Error(err))
		response.FailWithMessage("<UNK>", c)
		return
	}
	response.OkWithData(user, c)
}
