package global

import (
	"gorm.io/gorm"
	"shop_srvs/order_srv/proto"
)

var (
	DB                 *gorm.DB
	GoodsSrvClient     proto.GoodsClient
	InventorySrvClient proto.InventoryClient
	Tags               = []string{"linko", "lin", "order_srv"}
)

const (
	Host                 = "192.168.2.106" //数据库
	Port                 = "3306"
	User                 = "root"
	Password             = "root"
	ServerPort           = 8006 //自己启动的服务
	ServerHost           = "127.0.0.1"
	ServerConsulHost     = "192.168.2.106" //注册中心
	SeverConsulPort      = 8500
	ServerName           = "order_srv"
	InventorySrvInfoName = "inventory_srv" //连接的服务
	GoodsSrvInfoName     = "goods_srv"
	RocketMQ             = "192.168.2.106:9876"
)
