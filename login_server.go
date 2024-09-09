package main

import (
	"three-kingdoms-go/config"
	"three-kingdoms-go/net"
	"three-kingdoms-go/server/login"
)

// localhost:8888 服务 api/login 路由
//ws 区别 ws://localhost:8899 服务器 发消息 （封账为路由）

func LoginServer() {
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "8003")

	s := net.NewServer(host + ":" + port)
	login.Init()
	s.Router(login.Router)
	s.Start()
}
