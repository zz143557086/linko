package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_api/userop_web/proto"
)

var (
	Trans ut.Translator

	GoodsSrvClient proto.GoodsClient
	MessageClient  proto.MessageClient
	AddressClient  proto.AddressClient
	UserFavClient  proto.UserFavClient

	Tags = []string{"linko", "lin", "userop_web"}
)

const (
	ClientPort       = 8024
	ClientHost       = "127.0.0.1"
	ServerName       = "userop_srv"
	ClientName       = "userop_web"
	ConsulHost       = "192.168.2.106"
	ConsulPort       = 8500
	Key              = "ElwvtT$G@ceMdGE@rsBAZc9vOFs9zqO4"
	GoodsSrvInfoName = "goods_srv"
)
