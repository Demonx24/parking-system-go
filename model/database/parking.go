package database

import "time"

type ParkingLot struct { //停车场信息
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"size:128;not null" json:"name"`
	Address        string    `gorm:"size:256" json:"address"`
	Latitude       float64   `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude      float64   `gorm:"type:decimal(10,7)" json:"longitude"`
	TotalSlots     int       `gorm:"default:0" json:"total_slots"`
	AvailableSlots int       `gorm:"default:0" json:"available_slots"`
	PricePerHour   float64   `gorm:"type:decimal(10,2);default:0" json:"price_per_hour"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
type ParkingRecord struct { //停车记录表
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64     `gorm:"not null" json:"user_id"`
	ParkingLotID uint64     `gorm:"not null" json:"parking_lot_id"`
	PlateNumber  string     `gorm:"type:varchar(20);not null" json:"plate_number"`
	EntryTime    time.Time  `gorm:"not null" json:"entry_time"`
	ExitTime     *time.Time `json:"exit_time"`
	TotalFee     float64    `gorm:"type:decimal(10,2);default:0" json:"total_fee"`
	Status       int8       `gorm:"default:0" json:"status"` // 0未出场，1已出场
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// 车位状态表
type ParkingSlot struct {
	ID           uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ParkingLotID uint64    `gorm:"not null" json:"parking_lot_id"`
	SlotNumber   string    `gorm:"size:16" json:"slot_number"` //车位编号
	Status       int8      `gorm:"default:0" json:"status"`
	UpdatedAt    time.Time `json:"updated_at"`
}
