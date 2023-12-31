package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"net"
	"os"
	"os/signal"
	"shop_srvs/userop_srv/handler"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"shop_srvs/userop_srv/global"
	"shop_srvs/userop_srv/initialize"
	"shop_srvs/userop_srv/proto"
)

func main() {
	// 日志初始化
	initialize.InitLogger()
	// 数据库初始化
	initialize.InitDB()

	// 从命令行参数获取 IP 地址和端口号
	Ip := flag.String("ip", global.ServerHost, "IP地址")
	Port := flag.Int("port", global.ServerPort, "端口号")
	zap.S().Info("ip地址为", *Ip)
	zap.S().Info("端口号为", *Port)

	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册 gRPC 服务
	proto.RegisterAddressServer(server, &handler.UserOpServer{})
	proto.RegisterMessageServer(server, &handler.UserOpServer{})
	proto.RegisterUserFavServer(server, &handler.UserOpServer{})

	// 监听指定的 IP 地址和端口号
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *Ip, *Port))
	if err != nil {
		panic("监听失败: " + err.Error())
	}
	zap.S().Debugf("成功监听到")
	if global.DB == nil {
		zap.S().Debugf("指针为空")
	}
	// 在新的 Goroutine 中启动 gRPC 服务器
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("服务启动失败: " + err.Error())
		}
	}()

	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConsulHost, global.SeverConsulPort)

	// 创建 Consul 客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	zap.S().Debugf("注册服务")

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerName // 服务名称
	registration.ID = global.ServerName   // 服务ID
	registration.Port = *Port             // 服务端口号
	registration.Tags = global.Tags       // 标签，用于服务发现
	registration.Address = *Ip            // 服务 IP 地址
	//registration.Check = check
	//1. 如何启动两个服务
	//2. 即使我能够通过终端启动两个服务，但是注册到consul中的时候也会被覆盖

	// 注册服务到 Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic("服务注册失败: " + err.Error())
	}

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 注销服务
	if err = client.Agent().ServiceDeregister(global.ServerName); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")
}
