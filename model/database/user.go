package database

import "time"

type User struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	OpenID    string    `gorm:"size:64;unique;not null" json:"openid" form:"openid"`
	Nickname  string    `gorm:"size:64" json:"nickname" form:"nickname"`
	Phone     string    `gorm:"size:20" json:"phone" form:"phone"`
	CarPlate  string    `gorm:"size:16" json:"car_plate" form:"car_plate"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
