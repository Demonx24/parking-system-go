package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"parking-system-go/global"
	"parking-system-go/model/database"
	"parking-system-go/model/response"
)

type BarrierLogApi struct{}

func (b *BarrierLogApi) BarrierLogCreate(c *gin.Context) {
	var req database.BarrierLog
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("新建抬杆记录错误，前端参数错误", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if _, err := barrierLogService.Create(req); err != nil {
		global.Log.Error("新建抬杆记录错误，创建sql错误", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}
