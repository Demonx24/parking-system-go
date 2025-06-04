package main

import (
	"parking-system-go/core"
	"parking-system-go/global"
	"parking-system-go/initialize"
)

func main() {
	global.Config = core.InitConfig()
	global.DB = initialize.InitGorm()
	core.RunServer()
}
