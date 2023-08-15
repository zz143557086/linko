package main

import (
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"shop_srvs/inventory_srv/global"
	"shop_srvs/inventory_srv/handler"
	"shop_srvs/inventory_srv/initialize"
	"shop_srvs/inventory_srv/proto"
	"syscall"
)

func main() {
	// 日志初始化
	initialize.InitLogger()
	// 数据库初始化
	initialize.InitDB()
	//resdis锁的初始化
	initialize.InitRedis()
	// 从命令行参数获取 IP 地址和端口号
	Ip := flag.String("ip", global.ServerHost, "IP地址")
	Port := flag.Int("port", global.ServerPort, "端口号")
	zap.S().Info("ip地址为", *Ip)
	zap.S().Info("端口号为", *Port)

	// 创建 gRPC 服务器
	server := grpc.NewServer()

	// 注册 gRPC 服务
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})

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
	//监听归还topic

	// 创建一个新的PushConsumer实例，并指定NameServer地址和消费者组名
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{global.RocketMQ}),
		consumer.WithGroupName("shop_inventory"),
	)

	// 订阅主题为"order_reback"的消息，并指定消费函数
	if err := c.Subscribe("order_reback", consumer.MessageSelector{}, handler.AutoReback); err != nil {
		fmt.Println("读取消息失败")
	}

	_ = c.Start()
	//不能让主goroutine退出

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
