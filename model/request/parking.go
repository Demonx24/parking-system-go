package request

type ParkingStatus struct {
	CarPlate string `gorm:"size:16" json:"car_plate"`
	OpenID   string `gorm:"size:64;unique;not null" json:"openid"`
	Nickname string `gorm:"size:64" json:"nickname"`
	Phone    string `gorm:"size:20" json:"phone"`
}
