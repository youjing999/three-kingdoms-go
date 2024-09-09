package login

import (
	"three-kingdoms-go/db"
	"three-kingdoms-go/net"
	"three-kingdoms-go/server/login/controller"
)

var Router = net.NewRouter()

func Init() {
	//测试数据库，并且初始化
	db.TestDB()

	//还有其他初始化方法
	initRouter()
}
func initRouter() {
	controller.DefaultAccount.Router(Router)
}
