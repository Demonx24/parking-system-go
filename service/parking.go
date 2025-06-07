package service

import (
	"errors"
	"gorm.io/gorm"
	"parking-system-go/global"
	"parking-system-go/model/database"
)

type ParkingService struct{}

func (s *ParkingService) ParkingStatus(req database.ParkingRecord) (database.ParkingRecord, error) {
	var record database.ParkingRecord

	if req.ID == 0 && req.PlateNumber == "" {
		return record, errors.New("必须提供 ID 或车牌号")
	}

	query := global.DB.Model(&database.ParkingRecord{})
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	} else {
		query = query.Where("car_plate = ?", req.PlateNumber)
	}

	err := query.First(&record).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return record, err
	}

	return record, nil
}

func (parkingService *ParkingService) GetParkingLots(req database.ParkingLot) (database.ParkingLot, error) {
	var lot database.ParkingLot
	db := global.DB
	switch {
	case req.ID != 0:
		db = db.Where("id = ?", req.ID)
	case req.Name != "":
		db = db.Where("name = ?", req.Name)
	default:
		return database.ParkingLot{}, errors.New("必须提供 Id,停车场名字")
	}
	err := db.First(&lot).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return database.ParkingLot{}, err
	}
	return lot, nil
}
