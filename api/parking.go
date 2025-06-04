package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math"
	"parking-system-go/global"
	"parking-system-go/model/database"
	"parking-system-go/model/response"
	"time"
)

type ParkingApi struct{}

func (p *ParkingApi) GetParkingStatus(c *gin.Context) {
	var req database.User
	if err := c.ShouldBindQuery(&req); err != nil {
		global.Log.Error("查询停车状态参数错误", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}
	//根据前端的信息查询用户信息，提供车牌号去查询停车信息
	if req.CarPlate == "" {
		user, err := userService.GetUser(req)
		if err != nil {
			global.Log.Error("查询用户信息数据失败", zap.Error(err))
			response.FailWithMessage(err.Error(), c)
		}
		req = user
	}

	var parkingrecord = database.ParkingRecord{
		CarPlate: req.CarPlate,
	}
	parking, err := parkingService.ParkingStatus(parkingrecord)
	if err != nil {
		global.Log.Error("查询停车状态失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	}
	//封装响应数据
	var parkingStatus = response.ParkingStatus{
		CarPlate:  parking.CarPlate,
		TotalFee:  parking.TotalFee,
		EntryTime: parking.EntryTime,
		ExitTime:  parking.ExitTime,
	}
	//如果还未出场，将计算预计费用，入场到现在的费用
	if parkingStatus.ExitTime == nil {
		timeNow := time.Now()
		duration := timeNow.Sub(parkingStatus.EntryTime) // time.Duration，单位纳秒
		minutes := duration.Minutes()                    // float64，单位分钟
		//分钟为单位赋值给TotalTime
		parkingStatus.TotalTime = time.Duration(minutes)
		var lot = database.ParkingLot{
			ID: parking.ParkingLotID,
		}
		//查询停车场费用，返回给parking响应数据
		lot, err = parkingService.GetParkingLots(lot)
		if err != nil {
			global.Log.Error("查询停车场信息失败", zap.Error(err))
			response.FailWithMessage(err.Error(), c)
		}
		//保留三位小数
		factor := math.Pow(10, float64(3))
		parkingStatus.TotalFee = math.Round((lot.PricePerHour*minutes/60)*factor) / factor
	} else {
		parkingStatus.TotalTime = time.Duration(parkingStatus.ExitTime.Sub(parkingStatus.EntryTime).Minutes())
	}
	response.OkWithData(parkingStatus, c)
}
