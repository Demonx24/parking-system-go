package main

import (
	"fmt"
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
	global.WeChat = initialize.InitWeChat()
	core.RunServer()
	r := gin.New()
	r.Use(middleware.CustomRecovery(global.Log))
	fmt.Println(global.WeChat.WeChatAPIKey)

}
