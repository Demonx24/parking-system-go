package global

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"parking-system-go/config"
)

var (
	Config *config.Config
	DB     *gorm.DB
	Log    *zap.Logger
)
