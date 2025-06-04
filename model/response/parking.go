package response

import "time"

type ParkingStatus struct {
	CarPlate  string        `gorm:"size:16" json:"car_plate"`
	EntryTime time.Time     `gorm:"not null" json:"entry_time"`
	ExitTime  *time.Time    `json:"exit_time"`
	TotalTime time.Duration `gorm:"not null" json:"total_time"`                   //入场时间到当前时间的总共时间
	TotalFee  float64       `gorm:"type:decimal(3,2);default:0" json:"total_fee"` //总共费用
	Status    string        `json:"status"`                                       // 0未出场，1已出场，返回string到前端
}
