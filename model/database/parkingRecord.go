package database

import "time"

type ParkingRecord struct {
	ID           uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint64     `gorm:"not null" json:"user_id"`
	ParkingLotID uint64     `gorm:"not null" json:"parking_lot_id"`
	CarPlate     string     `gorm:"size:16" json:"car_plate"`
	EntryTime    time.Time  `gorm:"not null" json:"entry_time"`
	ExitTime     *time.Time `json:"exit_time"`
	TotalFee     float64    `gorm:"type:decimal(10,2);default:0" json:"total_fee"`
	Status       int8       `gorm:"default:0" json:"status"` // 0未出场，1已出场
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
