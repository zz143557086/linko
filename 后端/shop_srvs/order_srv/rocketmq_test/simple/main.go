package main

import (
	"context"
	"fmt"
	"shop_srvs/order_srv/global"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	// 创建一个新的Producer实例，并指定NameServer地址
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{global.RocketMQ}))
	if err != nil {
		panic("生成producer失败")
	}

	// 启动Producer
	if err = p.Start(); err != nil {
		panic("启动producer失败")
	}

	// 创建一条消息，并指定Topic和消息内容
	message := primitive.NewMessage("imooc1", []byte("this is imooc1"))

	// 同步发送消息，并接收响应结果
	res, err := p.SendSync(context.Background(), message)
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	} else {
		fmt.Printf("发送成功: %s\n", res.String())
	}
	time.Sleep(time.Hour)
	// 关闭Producer
	if err = p.Shutdown(); err != nil {
		panic("关闭producer失败")
	}
}
