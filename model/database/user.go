package database

import "time"

type User struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
	OpenID      string    `gorm:"column:openid;size:64;unique;not null" json:"openid" form:"openid"`
	Nickname    string    `gorm:"size:64" json:"nickname" form:"nickname"`
	Phone       string    `gorm:"size:20" json:"phone" form:"phone"`
	PlateNumber string    `gorm:"type:varchar(20);not null" json:"plate_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
