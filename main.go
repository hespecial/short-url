package main

import (
	"short-url/global"
	"short-url/internal/router"
)

func main() {
	global.InitConfig()
	global.InitLogger()
	global.InitMysql()
	global.InitRedis()

	router.StartServer()
}
