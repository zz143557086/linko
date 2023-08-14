package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"shop_srvs/goods_srv/proto"
)

var brandClient proto.GoodsClient
var conn *grpc.ClientConn

func TestGetCategoryBrandList() {
	rsp, err := brandClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{Pages: 5, PagePerNums: 10})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	fmt.Println(rsp.Data)
}

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:8005", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	brandClient = proto.NewGoodsClient(conn)
}

func main() {
	Init()
	//TestCreateUser()
	//TestGetCategoryList()
	TestGetCategoryBrandList()

	conn.Close()
}
