package service

import (
	"errors"
	"gorm.io/gorm"
	"parking-system-go/global"
	"parking-system-go/model/database"
)

type ParkingService struct{}

func (parkingService *ParkingService) ParkingStatus(req database.ParkingRecord) (database.ParkingRecord, error) {
	var record database.ParkingRecord
	db := global.DB
	switch {
	case req.ID != 0:
		db = db.Where("id = ?", req.ID)
	case req.CarPlate != "":
		db = db.Where("car_plate = ?", req.CarPlate)
	default:
		return database.ParkingRecord{}, errors.New("必须提供 Id,车牌号")
	}
	err := db.First(&record).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return database.ParkingRecord{}, err
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
