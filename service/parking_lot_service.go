package service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"parking-system-go/global"
	"parking-system-go/model/database"
	"parking-system-go/utils"
	"time"
)

type ParkingLotService struct{}

func (s *ParkingLotService) Create(lot *database.ParkingLot) error {
	return global.DB.Create(lot).Error
}

func (s *ParkingLotService) GetParkingLot(req database.ParkingLot) (database.ParkingLot, error) {
	var lot database.ParkingLot
	if req.ID == 0 && req.Name == "" {
		return lot, errors.New("必须提供 ID 或停车场名称")
	}

	query := global.DB.Model(&database.ParkingLot{})
	if req.ID != 0 {
		query = query.Where("id = ?", req.ID)
	} else {
		query = query.Where("name = ?", req.Name)
	}

	err := query.First(&lot).Error
	if err != nil {
		return lot, err
	}
	return lot, nil
}

func (s *ParkingLotService) Update(lot *database.ParkingLot) error {
	return global.DB.Save(lot).Error
}

func (s *ParkingLotService) Delete(where *database.ParkingLot) error {
	return global.DB.Where(where).Delete(&database.ParkingLot{}).Error
}

// 查询redis中每个停车场信息，如果redis中信息不存在就存入
func (s *ParkingLotService) GetParkingLotR(lot *database.ParkingLot) (database.ParkingLot, error) {
	fmt.Printf("global.Redis: %#v\n", global.Redis)
	if lot.ID == 0 && lot.Name == "" {
		return *lot, errors.New("必须提供 ID 或停车场名称")
	}
	key := fmt.Sprintf("parking_lot:%d", lot.ID)

	return utils.GetOrSetStruct(
		key,
		time.Hour,
		func() (database.ParkingLot, error) {
			var pl database.ParkingLot
			if lot.ID != 0 {
				err := global.DB.Where("id = ?", lot.ID).First(&pl).Error
				return pl, err
			} else {
				err := global.DB.Where("name = ?", lot.Name).First(&pl).Error
				return pl, err
			}
		},
	)
}

// 入场
func (s *ParkingLotService) DecrementAvailableSlotsWithPessimisticLock(parkingLotID uint64) error {
	// 开启事务
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var lot database.ParkingLot

		// 查询并锁定记录，FOR UPDATE 语法（悲观锁）
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", parkingLotID).
			First(&lot).Error; err != nil {
			return err
		}

		// 检查剩余车位
		if lot.AvailableSlots <= 0 {
			return errors.New("停车场剩余车位不足")
		}

		// 剩余车位减一
		lot.AvailableSlots -= 1

		// 保存修改
		if err := tx.Save(&lot).Error; err != nil {
			return err
		}

		return nil
	})
}

// 出场
func (s *ParkingLotService) IncrementAvailableSlotsWithPessimisticLock(parkingLotID uint64) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		var lot database.ParkingLot
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", parkingLotID).
			First(&lot).Error; err != nil {
			return err
		}

		// 可加最大车位数限制判断

		lot.AvailableSlots += 1

		if err := tx.Save(&lot).Error; err != nil {
			return err
		}
		return nil
	})
}
