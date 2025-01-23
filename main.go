package main

import (
	"short-url/api"
	"short-url/global"
)

func main() {
	global.InitConfig()
	global.InitLogger()
	global.InitMysql()
	global.InitRedis()

	api.StartServer()
}
