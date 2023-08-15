package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_api/userop_web/global"
	"shop_api/userop_web/proto"
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

	userOpSrvHost := ""
	userOpSrvPort := 0

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	data, err := client.Agent().ServicesWithFilter(fmt.Sprintf("Service == \"%s\"", global.ServerName))
	//data, err := client.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, global.ServerConfig.UserSrvInfo.Name))
	if err != nil {
		panic(err)
	}
	for _, value := range data {
		userOpSrvHost = value.Address
		userOpSrvPort = value.Port
		break
	}
	zap.S().Debugf("拨号ip为：%s:%d", userOpSrvHost, userOpSrvPort)
	if userOpSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【用户留言服务失败】")
		return
	}

	//拨号连接用户grpc服务器 跨域的问题 - 后端解决 也可以前端来解决
	userOpConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userOpSrvHost, userOpSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetuserOpList] 连接 【用户留言服务失败】",
			"msg", err.Error(),
		)
	}
	//1. 后续的用户服务下线了 2. 改端口了 3. 改ip了 负载均衡来做

	//2. 已经事先创立好了连接，这样后续就不用进行再次tcp的三次握手
	//3. 一个连接多个groutine共用，性能 - 连接池
	global.UserFavClient = proto.NewUserFavClient(userOpConn)
	global.MessageClient = proto.NewMessageClient(userOpConn)
	global.AddressClient = proto.NewAddressClient(userOpConn)
	// 使用 Zap 输出一条调试日志，表示连接服务成功
	zap.S().Debug("成功连接库存服务")
}
