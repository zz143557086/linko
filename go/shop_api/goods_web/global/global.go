package global

import (
	ut "github.com/go-playground/universal-translator"

	"shop_api/goods_web/proto"
)

var (
	Trans ut.Translator

	GoodsSrvClient proto.GoodsClient

	Tags = []string{"linko", "lin", "user_web"}
)

const (
	ClientPort = 8022            //启动服务的端口号
	ClientHost = "127.0.0.1"     //启动服务的ip
	ServerName = "goods_srv"     //远程调用服务的名称
	ClientName = "goods_web"     //启动服务的名称
	ConsulHost = "192.168.2.106" //注册中心的ip
	ConsulPort = 8500            //注册中心的端口号
	Key        = "ElwvtT$G@ceMdGE@rsBAZc9vOFs9zqO4"
)
