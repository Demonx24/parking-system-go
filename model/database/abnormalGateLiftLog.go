package database

import "time"

type AbnormalGateLiftLog struct {
	ID              uint64    `gorm:"primaryKey" json:"id"`
	ParkingRecordID *uint64   `json:"parking_record_id"` // 可为 nil
	CarPlate        string    `json:"car_plate"`
	ParkingLotID    uint64    `json:"parking_lot_id"`
	Operator        string    `json:"operator"`
	Reason          string    `json:"reason"`
	SnapshotURL     string    `json:"snapshot_url"`
	GateTime        time.Time `json:"gate_time"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
