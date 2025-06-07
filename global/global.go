package global

import (
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"parking-system-go/config"
)

var (
	Config *config.Config
	DB     *gorm.DB
	Log    *zap.Logger
	Redis  redis.Client
	WeChat *config.WeChat
)
