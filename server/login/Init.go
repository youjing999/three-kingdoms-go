package login

import (
	"three-kingdoms-go/net"
	"three-kingdoms-go/server/login/controller"
)

var Router = net.NewRouter()

func Init() {
	//还有其他初始化方法
	initRouter()
}
func initRouter() {
	controller.DefaultAccount.Router(Router)
}
