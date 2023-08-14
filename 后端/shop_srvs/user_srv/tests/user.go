package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"shop_srvs/user_srv/proto"
)

var userClient proto.UserClient
var conn *grpc.ClientConn

func Init() {
	// 初始化连接 gRPC 服务器
	var err error
	conn, err = grpc.Dial("127.0.0.1:8003", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	userClient = proto.NewUserClient(conn)
}

func TestGetUserList() {
	// 获取用户列表
	rsp, err := userClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    1,
		PSize: 5,
	})
	if err != nil {
		panic(err)
	}
	for _, user := range rsp.Data {
		// 打印用户信息
		fmt.Println(user)

		if err != nil {
			panic(err)
		}
	}
}

func TestCreateUser() {
	// 创建用户+
	for i := 0; i < 10; i++ {
		rsp, err := userClient.CreateUser(context.Background(), &proto.CreateUserInfo{
			Name:     fmt.Sprintf("joker%d", i),
			Mobile:   fmt.Sprintf("1111111111%d", i),
			PassWord: "123456",
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(rsp.Id)
	}
}

func TestGetUserByMobile(mobile string) {
	// 通过手机号获取用户信息
	rsp, err := userClient.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: mobile})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp)
}

func main() {
	// 初始化连接
	Init()
	//TestCreateUser()
	// 测试获取用户列表
	TestGetUserList()
	// 根据手机号获取用户信息
	TestGetUserByMobile("15520050430")
	// 关闭连接
	conn.Close()
}
