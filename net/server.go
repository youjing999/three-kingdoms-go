package net

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Server struct {
	addr   string
	router *Router
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) Router(router *Router) {
	s.router = router
}

// Start 启动服务
func (s *Server) Start() {
	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}

}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. http协议升级成websocket
	wsConn, err := wsUpgrade.Upgrade(w, r, nil)
	if err != nil {
		//打印日志 同时 退出应用程序
		log.Println("websocket服务连接出错", err)
	}
	//websocket通道建立之后 不管是客户端还是服务端 都可以收发消息
	//发消息的时候 把消息当做路由 来去处理 消息是有格式的，先定义消息的格式
	//客户端 发消息的时候 {Name:"account.login"} 收到之后 进行解析 认为想要处理登录逻辑
	wsServer := NewWsServer(wsConn)
	wsServer.Router(s.router)
	wsServer.Start()
	wsServer.Handshake()
}

// http升级websocket协议的配置
var wsUpgrade = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
