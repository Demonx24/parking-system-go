package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math"
	"parking-system-go/global"
	"parking-system-go/model/database"
	"parking-system-go/model/response"
	"parking-system-go/utils"
	"time"
)

type ParkingApi struct{} // //////////////////////////////////////////////////////////////////////////
// 入场信息
func (p *ParkingApi) Entry(c *gin.Context) {
	//前端传过来抬杆信息
	var req database.BarrierLog
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("抬杆信息失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	inParking, err := barrierLogService.IsCarInParking(req.PlateNumber)
	if err != nil {
		response.FailWithMessage("判断车辆是否在场失败", c)
		return
	}
	if inParking {
		response.OkWithMessage("车辆在停车场内", c)
		return
	}

	req.LaneType = "entry"
	timestamp := time.Now()
	req.Timestamp = timestamp
	//将抬杆信息保存
	barrierLog, err := barrierLogService.Create(req)
	if err != nil {
		global.Log.Error("保存抬杆信息失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	//停车场剩余车位

	if err := parkinglotService.DecrementAvailableSlotsWithPessimisticLock(*barrierLog.ParkingID); err != nil {
		global.Log.Error("修改当前停车场剩余车位失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	parkingRecord := database.ParkingRecord{
		ParkingLotID: *barrierLog.ParkingID,
		PlateNumber:  barrierLog.PlateNumber,
		Status:       0,
		EntryTime:    barrierLog.Timestamp,
		CreatedAt:    barrierLog.CreatedAt,
	}
	if err := parkingrecordService.Create(&parkingRecord); err != nil {
		global.Log.Error("保存停车信息失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("欢迎"+parkingRecord.PlateNumber, c)
}

// 出场信息
func (p *ParkingApi) Exit(c *gin.Context) {
	var req database.BarrierLog
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("出场参数绑定失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	parkingRecord := database.ParkingRecord{
		PlateNumber: req.PlateNumber,
	}
	// 查找未出场的停车记录
	parkingRecord, err := parkingrecordService.GetRecord(parkingRecord)
	if err != nil {
		global.Log.Error("未找到在场车辆记录", zap.Error(err))
		response.FailWithMessage("未找到该车辆在场记录", c)
		return
	}
	if parkingRecord.Status != 0 {
		response.FailWithMessage("车辆未进入停车场", c)
		return
	}
	// 保存出场抬杆记录
	req.LaneType = "exit"
	req.Timestamp = time.Now()
	_, err = barrierLogService.Create(req)
	if err != nil {
		global.Log.Error("保存出场抬杆记录失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 更新出场时间
	parkingRecord.ExitTime = &req.Timestamp
	parkingRecord.Status = 0 // 已出场
	//算停车费用
	ps, err := buildParkingStatus(parkingRecord)
	if err != nil {
		global.Log.Error("解算停车费失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	parkingRecord.TotalFee = ps.TotalFee
	order := database.Order{
		Amount:          ps.TotalFee,
		PlateNumber:     parkingRecord.PlateNumber,
		Status:          0,
		ParkingRecordID: parkingRecord.ID,
		UserID:          parkingRecord.UserID,
		OrderID:         utils.GenerateOrderID(),
	}
	// 保存更新的停车记录
	if err := parkingrecordService.Update(&parkingRecord); err != nil {
		global.Log.Error("更新出场信息失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	// ==== 模拟支付逻辑（可扩展） ====
	var wxResp response.UnifiedOrderResponse
	wxResp.CodeURL, wxResp.PrepayID, err = payWeChatService.CreatePayment(&order, c.ClientIP())
	if err != nil {
		global.Log.Error("调用微信统一接口失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(wxResp, c)

}

// 查询停车状态
func (p *ParkingApi) GetParkingStatus(c *gin.Context) {
	var req database.User
	if err := c.ShouldBindQuery(&req); err != nil {
		global.Log.Error("查询停车状态参数错误", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 如果前端没有传车牌，则用其他信息查询用户，补充车牌
	if req.PlateNumber == "" {
		user, err := userService.GetUser(req)
		if err != nil {
			global.Log.Error("查询用户信息数据失败", zap.Error(err))
			response.FailWithMessage(err.Error(), c)
			return
		}
		req.PlateNumber = user.PlateNumber
	}

	parking, err := parkingService.ParkingStatus(database.ParkingRecord{PlateNumber: req.PlateNumber})
	if err != nil {
		global.Log.Error("查询停车状态失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	parkingStatus, err := buildParkingStatus(parking)
	if err != nil {
		global.Log.Error("计算停车状态失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(parkingStatus, c)
}

// 抽离计算停车状态和费用逻辑
func buildParkingStatus(parking database.ParkingRecord) (response.ParkingStatus, error) {

	ps := response.ParkingStatus{
		CarPlate:  parking.PlateNumber,
		EntryTime: parking.EntryTime,
		ExitTime:  parking.ExitTime,
	}

	var minutes float64
	if parking.ExitTime == nil {
		minutes = time.Since(parking.EntryTime).Minutes()
	} else {
		minutes = parking.ExitTime.Sub(parking.EntryTime).Minutes()
	}

	// 转换为 duration（分钟数）
	ps.TotalTime = time.Duration(minutes) * time.Minute

	// 查询停车场价格
	lot, err := parkingService.GetParkingLots(database.ParkingLot{ID: parking.ParkingLotID})
	if err != nil {
		return ps, err
	}

	// 计算费用，保留3位小数
	factor := math.Pow(10, 3)
	ps.TotalFee = math.Round((lot.PricePerHour*minutes/60)*factor) / factor

	return ps, nil
}

// 查询所有停车场信息
func (p *ParkingApi) GetParkingLotR(c *gin.Context) {
	var req database.ParkingLot
	if err := c.ShouldBindJSON(&req); err != nil {
		global.Log.Error("查询参数错误", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	req, err := parkinglotService.GetParkingLotR(&req)
	if err != nil {
		global.Log.Error("查询信息失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(req, c)
}
