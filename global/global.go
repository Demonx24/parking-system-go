package global

import (
	"gorm.io/gorm"
	"parking-system-go/config"
)

var (
	Config *config.Config
	DB     *gorm.DB
)
