package controller

import (
	"github.com/mitchellh/mapstructure"
	"log"
	"three-kingdoms-go/constant"
	"three-kingdoms-go/db"
	"three-kingdoms-go/net"
	"three-kingdoms-go/server/login/model"
	"three-kingdoms-go/server/login/proto"
	"three-kingdoms-go/server/models"
	"three-kingdoms-go/utils"
	"time"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(r *net.Router) {
	group := r.Group("account")
	group.AddRouter("login", a.login)

}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {

	/**
	1.用户名 密码 硬件Id
	2.根据用户名 查询user表 得到数据
	3.进行密码比对 登录成功
	4.保存用户的登录记录
	5.保存用户的最后一次登录信息
	6.客户端 要一个session jwt生成一个加密字符串的算法
	7. 客户端 在发起要用户灯枯的行为时，判断用户是否合法
	*/
	loginReq := &proto.LoginReq{}
	loginRes := &proto.LoginRsp{}
	mapstructure.Decode(req.Body.Msg, loginReq)
	user := &models.User{}
	ok, err := db.Engine.Table(user).Where("username=?", loginReq.Username).Get(user)
	if err != nil {
		log.Println("用户表查询出错", err)
		return
	}
	if !ok {
		//有没有查出来数据
		rsp.Body.Code = constant.UserNotExist
		return
	}
	pwd := utils.Password(loginReq.Password, user.Passcode)
	if pwd != user.Passwd {
		rsp.Body.Code = constant.PwdIncorrect
		return
	}
	//jwt A.B.C 三部分 A定义加密算法 B定义放入的数据 C部分 根据秘钥+A和B生成加密字符串
	token, _ := utils.Award(user.UId)
	rsp.Body.Code = constant.OK
	loginRes.UId = user.UId
	loginRes.Username = user.Username
	loginRes.Session = token
	loginRes.Password = ""
	rsp.Body.Msg = loginRes

	//保存用户登录记录
	ul := &model.LoginHistory{
		UId: user.UId, CTime: time.Now(), Ip: loginReq.Ip,
		Hardware: loginReq.Hardware, State: model.Login,
	}
	db.Engine.Table(ul).Insert(ul)
	//最后一次登录的状态记录
	ll := &model.LoginLast{}
	ok, _ = db.Engine.Table(ll).Where("uid=?", user.UId).Get(ll)
	if ok {
		//有数据 更新
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		db.Engine.Table(ll).Update(ll)
	} else {
		ll.IsLogout = 0
		ll.Ip = loginReq.Ip
		ll.LoginTime = time.Now()
		ll.Session = token
		ll.Hardware = loginReq.Hardware
		ll.UId = user.UId
		_, err := db.Engine.Table(ll).Insert(ll)
		if err != nil {
			log.Println(err)
		}
	}

}
