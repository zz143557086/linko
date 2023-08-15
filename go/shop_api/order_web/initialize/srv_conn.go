package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_api/order_web/global"
	"shop_api/order_web/proto"
)

func InitSrvConn1() {
	//从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ConsulHost, global.ConsulPort)

	GoodsSrvHost := ""
	GoodsSrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.GoodsSrvInfoName))
	//data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		GoodsSrvHost = value.Address
		GoodsSrvPort = value.Port
		break
	}
	zap.S().Debugf("拨号ip为：%s:%d", GoodsSrvHost, GoodsSrvPort)
	if GoodsSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【商品服务失败】")
		return
	}

	GoodsConn, err := grpc.Dial(fmt.Sprintf("%s:%d", GoodsSrvHost, GoodsSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetGoodsList] 连接 【商品服务失败】",
			"msg", err.Error(),
		)
	}
	//1. 后续的用户服务下线了 2. 改端口了 3. 改ip了 负载均衡来做

	//2. 已经事先创立好了连接，这样后续就不用进行再次tcp的三次握手
	//3. 一个连接多个groutine共用，性能 - 连接池
	GoodsSrvClient := proto.NewGoodsClient(GoodsConn)
	global.GoodsSrvClient = GoodsSrvClient
	// 使用 Zap 输出一条调试日志，表示连接服务成功
	zap.S().Debug("成功连接商品服务")
}

func InitSrvConn2() {
	//从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ConsulHost, global.ConsulPort)

	inventorySrvHost := ""
	InventorySrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.InventorySrvInfoName))
	//data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		inventorySrvHost = value.Address
		InventorySrvPort = value.Port
		break
	}
	zap.S().Debugf("拨号ip为：%s:%d", inventorySrvHost, InventorySrvPort)
	if inventorySrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【库存服务失败】")
		return
	}

	//拨号连接用户grpc服务器 跨域的问题 - 后端解决 也可以前端来解决
	inventoryConn, err := grpc.Dial(fmt.Sprintf("%s:%d", inventorySrvHost, InventorySrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetInventoryList] 连接 【库存服务失败】",
			"msg", err.Error(),
		)
	}
	//1. 后续的用户服务下线了 2. 改端口了 3. 改ip了 负载均衡来做

	//2. 已经事先创立好了连接，这样后续就不用进行再次tcp的三次握手
	//3. 一个连接多个groutine共用，性能 - 连接池
	InventorySrvClient := proto.NewInventoryClient(inventoryConn)
	global.InventorySrvClient = InventorySrvClient
	// 使用 Zap 输出一条调试日志，表示连接服务成功
	zap.S().Debug("成功连接库存服务")
}
func InitSrvConn3() {
	//从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ConsulHost, global.ConsulPort)

	OrderSrvHost := ""
	OrderSrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.OrderSrvInfoName))
	//data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		OrderSrvHost = value.Address
		OrderSrvPort = value.Port
		break
	}
	zap.S().Debugf("拨号ip为：%s:%d", OrderSrvHost, OrderSrvPort)
	if OrderSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【订单服务失败】")
		return
	}

	//拨号连接用户grpc服务器 跨域的问题 - 后端解决 也可以前端来解决
	OrderConn, err := grpc.Dial(fmt.Sprintf("%s:%d", OrderSrvHost, OrderSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetOrderList] 连接 【订单服务失败】",
			"msg", err.Error(),
		)
	}
	//1. 后续的用户服务下线了 2. 改端口了 3. 改ip了 负载均衡来做

	//2. 已经事先创立好了连接，这样后续就不用进行再次tcp的三次握手
	//3. 一个连接多个groutine共用，性能 - 连接池
	OrderSrvClient := proto.NewOrderClient(OrderConn)
	global.OrderSrvClient = OrderSrvClient
	// 使用 Zap 输出一条调试日志，表示连接服务成功
	zap.S().Debug("成功连接订单服务")
}
