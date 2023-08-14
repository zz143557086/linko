package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_api/goods_web/global"
	"shop_api/goods_web/proto"
)

func InitSrvConn2() {
	//从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ConsulHost, global.ConsulPort)

	goodsSrvHost := ""
	goodsSrvPort := 0

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
		goodsSrvHost = value.Address
		goodsSrvPort = value.Port
		break
	}
	zap.S().Debugf("拨号ip为：%s:%d", goodsSrvHost, goodsSrvPort)
	if goodsSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}

	//拨号连接用户grpc服务器 跨域的问题 - 后端解决 也可以前端来解决
	goodsConn, err := grpc.Dial(fmt.Sprintf("%s:%d", goodsSrvHost, goodsSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}
	//1. 后续的用户服务下线了 2. 改端口了 3. 改ip了 负载均衡来做

	//2. 已经事先创立好了连接，这样后续就不用进行再次tcp的三次握手
	//3. 一个连接多个groutine共用，性能 - 连接池
	goodsSrvClient := proto.NewGoodsClient(goodsConn)
	global.GoodsSrvClient = goodsSrvClient
	// 使用 Zap 输出一条调试日志，表示连接服务成功
	zap.S().Debug("成功连接服务")
}
