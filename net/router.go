package net

import (
	"strings"
)

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)

// account login/loginOut
type group struct {
	prefix     string
	handlerMap map[string]HandlerFunc
}
type Router struct {
	group []*group
}

func NewRouter() *Router {
	return &Router{}
}

func (g *group) AddRouter(name string, handlerFunc HandlerFunc) {
	g.handlerMap[name] = handlerFunc
}

func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix:     prefix,
		handlerMap: make(map[string]HandlerFunc),
	}
	r.group = append(r.group, g)
	return g
}

func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	h := g.handlerMap[name]
	if h != nil {
		h(req, rsp)
	}
}
func (r *Router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	//req.Body.Name 路径 登录业务 account.login （account组标识）login 路由标识
	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""

	if len(strs) == 2 {
		prefix = strs[0]
		name = strs[1]
	}

	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, rsp)
			//fmt.Println(name)
		}
	}
}
