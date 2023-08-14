package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"shop_api/oss_web/global"
	"shop_api/oss_web/initialize"
	"syscall"
)

func main() {
	//1. 初始化logger
	initialize.InitLogger()

	//3. 初始化routers
	Router := initialize.Routers()
	//4. 初始化翻译
	if err := initialize.InitTrans("zh"); err != nil {
		panic(err)
	}

	// 服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ConsulHost, global.ConsulPort)
	// 创建 Consul 客户端
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	zap.S().Debugf("注册服务")

	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ClientName    // 服务名称
	registration.ID = global.ClientName      // 服务ID
	registration.Port = global.ClientPort    // 服务端口号
	registration.Tags = global.Tags          // 标签，用于服务发现
	registration.Address = global.ClientHost // 服务 IP 地址
	//registration.Check = check
	//1. 如何启动两个服务
	//2. 即使我能够通过终端启动两个服务，但是注册到consul中的时候也会被覆盖

	// 注册服务到 Consul
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic("服务注册失败: " + err.Error())
	}

	go func() {
		// 启动服务器并监听指定端口号
		if err := Router.Run(fmt.Sprintf(":%d", global.ClientPort)); err != nil {
			zap.S().Panicf("启动失败: %s", err.Error())
		}
	}()

	// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// 注销服务
	if err = client.Agent().ServiceDeregister(global.ClientName); err != nil {
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")

}
