package service

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"parking-system-go/global"
	"parking-system-go/model/database"
)

type BarrierLogService struct{}

// 创建抬杆记录
func (s *BarrierLogService) Create(log database.BarrierLog) (database.BarrierLog, error) {
	if err := global.DB.Create(&log).Error; err != nil {
		global.Log.Error("创建抬杆记录", zap.Error(err))
		return log, err
	}
	return log, nil
}

// 根据对象获取抬杆记录
func (s *BarrierLogService) GetBarrierLog(req database.BarrierLog) (database.BarrierLog, error) {
	var log database.BarrierLog
	if req.ID == 0 && req.BarrierID == nil {
		return log, errors.New("必须提供 ID 或设备编号")
	}

	query := global.DB.Model(&database.BarrierLog{})
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	} else {
		query = query.Where("barrier_id = ?", req.BarrierID)
	}

	err := query.First(&log).Error
	if err != nil {
		return log, err
	}
	return log, nil
}

// 更新抬杆记录（全部字段更新）
func (s *BarrierLogService) Update(log database.BarrierLog) (database.BarrierLog, error) {
	// 这里用Save，会根据主键更新全部字段
	if err := global.DB.Save(&log).Error; err != nil {
		global.Log.Error("更新抬杆记录（全部字段更新）", zap.Error(err))
		return log, err
	}
	return log, nil
}

// 删除抬杆记录
func (s *BarrierLogService) Delete(id uint64) error {
	if err := global.DB.Delete(&database.BarrierLog{}, id).Error; err != nil {
		global.Log.Error("删除抬杆记录", zap.Error(err))
		return err
	}
	return nil
}

// 查询全部抬杆记录（示例不带分页）
func (s *BarrierLogService) ListAll() ([]database.BarrierLog, error) {
	var logs []database.BarrierLog
	if err := global.DB.Order("id DESC").Find(logs).Error; err != nil {
		global.Log.Error("查询全部抬杆记录（示例不带分页）", zap.Error(err))
		return nil, err
	}
	return logs, nil
}

// 判断某车牌是否在场内（最后一次 entry > exit 则说明在场）
func (s *BarrierLogService) IsCarInParking(plateNumber string) (bool, error) {
	var lastEntry, lastExit database.BarrierLog

	// 获取最后一次 entry
	err := global.DB.
		Where("plate_number = ? AND lane_type = ?", plateNumber, "entry").
		Order("timestamp DESC").
		First(&lastEntry).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	// 获取最后一次 exit
	err = global.DB.
		Where("plate_number = ? AND lane_type = ?", plateNumber, "exit").
		Order("timestamp DESC").
		First(&lastExit).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	// 如果没有 entry，肯定不在场
	if lastEntry.ID == 0 {
		return false, nil
	}

	// 没有 exit，说明还没出场
	if lastExit.ID == 0 {
		return true, nil
	}

	// 比较时间
	return lastEntry.Timestamp.After(lastExit.Timestamp), nil
}
