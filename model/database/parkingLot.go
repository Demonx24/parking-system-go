package database

import "time"

type ParkingLot struct {
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
