package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_api/user_web/proto"
)

var (
	Trans         ut.Translator
	UserSrvClient proto.UserClient
	Tags          = []string{"linko", "lin", "user_web"}
)

const (
	ClientPort = 8021 //启动的端口号
	ServerPort = 8003
	Host       = "127.0.0.1"                        //数据库的ip
	ClientHost = "127.0.0.1"                        //启动的ip
	Key        = "ElwvtT$G@ceMdGE@rsBAZc9vOFs9zqO4" //jwt签名的key值
	ServerName = "user_srv"                         //远程调用grpc服务的名称
	ClientName = "user_web"                         //该服务的名称注册到consul注册中心
	ConsulHost = "192.168.2.106"                    //consul客服端的ip
	ConsulPort = 8500                               //consul客服端的端口号
)
