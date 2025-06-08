//go:build linux
// +build linux

package core

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// initServer 初始化一个 HTTP 服务器（适用于 Linux）
func initServer(address string, router *gin.Engine) server {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    10 * time.Minute,
		WriteTimeout:   10 * time.Minute,
		MaxHeaderBytes: 1 << 20,
	}
}
