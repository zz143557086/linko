package initialize

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"shop_api/user_web/global"
	"shop_api/user_web/proto"
)

func InitSrvConn() {
	// 声明错误变量
	var err error

	// 通过 gRPC Dial 函数与用户服务建立连接
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", global.Host, global.ServerPort), grpc.WithInsecure())
	if err != nil {
		// 如果连接失败，则使用 Zap 记录错误日志
		zap.S().Errorw("[GetUserList]连接[用户服务失败]", "msg", err.Error())
	}

	// 创建用户服务的 gRPC 客户端
	UserSrvClient := proto.NewUserClient(conn)

	// 将用户服务的 gRPC 客户端设置为全局变量
	global.UserSrvClient = UserSrvClient

	// 使用 Zap 输出一条调试日志，表示连接服务成功
	zap.S().Debug("成功连接服务")
}
func InitSrvConn2() {
	//从注册中心获取到用户服务的信息
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ConsulHost, global.ConsulPort)

	userSrvHost := ""
	userSrvPort := 0

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
		userSrvHost = value.Address
		userSrvPort = value.Port
		break
	}
	zap.S().Debugf("拨号ip为：%s:%d", userSrvHost, userSrvPort)
	if userSrvHost == "" {
		zap.S().Fatal("[InitSrvConn] 连接 【用户服务失败】")
		return
	}

	//拨号连接用户grpc服务器 跨域的问题 - 后端解决 也可以前端来解决
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", userSrvHost, userSrvPort), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【用户服务失败】",
			"msg", err.Error(),
		)
	}
	//1. 后续的用户服务下线了 2. 改端口了 3. 改ip了 负载均衡来做

	//2. 已经事先创立好了连接，这样后续就不用进行再次tcp的三次握手
	//3. 一个连接多个groutine共用，性能 - 连接池
	userSrvClient := proto.NewUserClient(userConn)
	global.UserSrvClient = userSrvClient
	// 使用 Zap 输出一条调试日志，表示连接服务成功
	zap.S().Debug("成功连接服务")
}
