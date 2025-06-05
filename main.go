package main

import (
	"github.com/gin-gonic/gin"
	"parking-system-go/core"
	"parking-system-go/global"
	"parking-system-go/initialize"
	"parking-system-go/middleware"
)

func main() {
	global.Config = core.InitConfig()
	global.DB = initialize.InitGorm()
	global.Log = core.InitLogger()
	core.RunServer()
	r := gin.New()
	r.Use(middleware.CustomRecovery(global.Log))
}
