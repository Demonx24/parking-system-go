package service

import (
	"errors"
	"parking-system-go/global"
	"parking-system-go/model/database"
)

type ParkingRecordService struct{}

func (s *ParkingRecordService) Create(record *database.ParkingRecord) error {
	return global.DB.Create(record).Error
}

func (s *ParkingRecordService) GetRecord(req database.ParkingRecord) (database.ParkingRecord, error) {
	var record database.ParkingRecord
	if req.ID == 0 && req.PlateNumber == "" {
		return record, errors.New("必须提供 ID 或车牌号")
	}

	query := global.DB.Model(&database.ParkingRecord{})
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	} else {
		query = query.Where("plate_number = ?", req.PlateNumber)
	}
	query = query.Order("created_at DESC")
	err := query.First(&record).Error
	if err != nil {
		return record, err
	}
	return record, nil
}

func (s *ParkingRecordService) Update(record *database.ParkingRecord) error {
	return global.DB.Save(record).Error
}

func (s *ParkingRecordService) Delete(where *database.ParkingRecord) error {
	return global.DB.Where(where).Delete(&database.ParkingRecord{}).Error
}
