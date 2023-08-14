package global

import (
	ut "github.com/go-playground/universal-translator"
	"shop_api/order_web/proto"
)

var (
	Trans ut.Translator

	GoodsSrvClient proto.GoodsClient

	OrderSrvClient proto.OrderClient

	InventorySrvClient proto.InventoryClient

	Tags = []string{"linko", "lin", "order_web"}
)

const (
	ClientPort           = 8023                               //启动服务的端口号
	ClientHost           = "127.0.0.1"                        //启动服务的ip
	ClientName           = "order_web"                        //启动服务 服务名称
	ConsulHost           = "192.168.2.106"                    //consul注册中心的ip
	ConsulPort           = 8500                               //consul注册中心的端口号
	Key                  = "ElwvtT$G@ceMdGE@rsBAZc9vOFs9zqO4" //jwt的签名密匙
	GoodsSrvInfoName     = "goods_srv"                        //远程调用的服务名称
	InventorySrvInfoName = "inventory_srv"
	OrderSrvInfoName     = "order_srv"
)
