package main

import (
	"context"
	"fmt"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

func main() {
	// 创建一个新的Producer实例，并指定NameServer地址为"192.168.0.104:9876"
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.2.106:9876"}))
	if err != nil {
		panic("生成producer失败")
	}

	// 启动Producer
	if err = p.Start(); err != nil {
		panic("启动producer失败")
	}

	// 创建一条延迟消息，并指定Topic和消息内容
	msg := primitive.NewMessage("imooc1", []byte("this is delay message"))
	msg.WithDelayTimeLevel(3) // 设置延迟级别为3，即延迟10s执行

	// 同步发送消息，并接收响应结果
	res, err := p.SendSync(context.Background(), msg)
	if err != nil {
		fmt.Printf("发送失败: %s\n", err)
	} else {
		fmt.Printf("发送成功: %s\n", res.String())
	}

	// 关闭Producer
	if err = p.Shutdown(); err != nil {
		panic("关闭producer失败")
	}

	// 支付的时候，淘宝、12306购票等场景，超时归还 - 定时执行逻辑
	// 可以使用延迟消息来实现，比如在下订单时设置一个合适的延迟时间
	// 当延迟时间到达时，执行特定的超时逻辑，同时消息中携带订单编号等信息方便查询

}
