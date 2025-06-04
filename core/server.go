package core

import (
	"fmt"
	"parking-system-go/global"
	"parking-system-go/initialize"
)

type server interface {
	ListenAndServe() error
}

// RunServer 用于启动服务器
func RunServer() {
	addr := global.Config.System.Addr()
	Router := initialize.InitRouter()

	// 加载所有的 JWT 黑名单，存入本地缓存
	//service.LoadAll()

	// 初始化服务器并启动
	s := initServer(addr, Router)
	fmt.Println(s.ListenAndServe().Error())
}
