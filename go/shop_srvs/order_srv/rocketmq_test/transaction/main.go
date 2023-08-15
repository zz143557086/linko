package main

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

// 自定义事务监听器
type OrderListener struct{}

// 执行本地事务逻辑
func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	fmt.Println("开始执行本地逻辑")
	time.Sleep(time.Second * 3)
	fmt.Println("执行本地逻辑失败")
	// 本地执行逻辑无缘无故失败，代码异常或宕机
	//return primitive.CommitMessageState 执行成功
	return primitive.CommitMessageState
}

// 回查本地事务状态
func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("rocketmq的消息回查")
	time.Sleep(time.Second * 15)
	return primitive.CommitMessageState
}

func main() {
	// 创建一个新的事务型Producer实例
	p, err := rocketmq.NewTransactionProducer(
		&OrderListener{}, // 使用自定义事务监听器
		producer.WithNameServer([]string{"192.168.2.106:9876"}), // 指定NameServer地址
	)
	if err != nil {
		panic("生成producer失败")
	}

	// 启动Producer
	if err = p.Start(); err != nil {
		panic("启动producer失败")
	}

	// 发送事务消息，并接收响应结果
	res, err := p.SendMessageInTransaction(context.Background(), primitive.NewMessage("order_reback", []byte("this is transaction message")))
	if err != nil {
		// 发送失败时打印错误消息
		fmt.Printf("发送失败: %s\n", err)
	} else {
		// 发送成功时打印响应结果
		fmt.Printf("发送成功: %s\n", res.String())
	}
	// 保持主goroutine运行一小时，以保证事务消息的回查和执行
	time.Sleep(time.Hour)

	// 关闭Producer
	if err = p.Shutdown(); err != nil {
		panic("关闭producer失败")
	}
}
