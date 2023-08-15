package main

import (
	"context"
	"fmt"
	"shop_srvs/order_srv/global"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
)

func main() {
	// 创建一个新的PushConsumer实例，并指定NameServer地址和消费者组名
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{global.RocketMQ}),
		consumer.WithGroupName("shop"),
	)

	// 订阅主题为"imooc1"的消息，并指定消费函数
	if err := c.Subscribe("imooc1", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i := range msgs {
			fmt.Printf("获取到值： %v \n", msgs[i])
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		fmt.Println("读取消息失败")
	}

	// 启动PushConsumer
	_ = c.Start()

	// 保持主goroutine运行一小时，以保证消息消费持续进行
	time.Sleep(time.Hour)

	// 关闭PushConsumer
	_ = c.Shutdown()
}
