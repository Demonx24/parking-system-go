package database

import "time"

type Order struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	OrderID         string     `gorm:"not null"  json:"order_id" form:"order_id"`
	UserID          uint64     `gorm:"not null" json:"user_id"`
	PlateNumber     string     `gorm:"type:varchar(20);not null" json:"plate_number"`
	ParkingRecordID uint64     `gorm:"not null" json:"parking_record_id"`
	Amount          float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status          int8       `gorm:"default:0" json:"status"` // 0待支付，1已支付，2失败
	PayTime         *time.Time `json:"pay_time"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	PrepayID        string     `gorm:"not null" json:"prepay_id" form:"prepay_id"`
}
