package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 用于吧panic也记录到zap中
func CustomRecovery(logger *zap.Logger) gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		// 使用 zap 记录 panic
		logger.Error("panic recovered",
			zap.Any("error", err),
			zap.String("url", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
		)

		// 返回 500 错误
		c.JSON(500, gin.H{"msg": "internal server error"})
		c.Abort()
	})
}
