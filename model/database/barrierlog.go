package database

import "time"

type BarrierLog struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	PlateNumber string    `gorm:"type:varchar(20);not null" json:"plate_number"`
	Timestamp   time.Time `gorm:"not null" json:"timestamp"`
	LaneType    string    `gorm:"type:enum('entry','exit');not null" json:"lane_type"`
	BarrierID   *string   `gorm:"type:varchar(64)" json:"barrier_id,omitempty"`
	ParkingID   *uint64   `json:"parking_id,omitempty"`
	Source      string    `gorm:"type:enum('auto','manual');default:'auto'" json:"source"`
	Result      string    `gorm:"type:enum('success','fail');default:'success'" json:"result"`
	Message     *string   `gorm:"type:varchar(255)" json:"message,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
