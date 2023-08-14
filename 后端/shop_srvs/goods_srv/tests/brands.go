package main

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"shop_srvs/goods_srv/proto"
)

var brandClient proto.GoodsClient
var conn *grpc.ClientConn

func TestGetBrandList() {
	rsp, err := brandClient.BrandList(context.Background(), &proto.BrandFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}
func TestCreateBrandList() {
	rsp, err := brandClient.CreateBrand(context.Background(), &proto.BrandRequest{
		Name: "世界",
		Logo: "按时间到了就爱上",
	})
	if err != nil {
		panic(err.Error() + "错误")
	}
	fmt.Println(rsp.Id)

}
func deleteBrand() {
	rsp, err := brandClient.DeleteBrand(context.Background(), &proto.BrandRequest{
		Id: 1111,
	})
	if err != nil {
		panic(err.Error() + "已经删除或者失败")

	}
	fmt.Println(rsp, "删除成功")
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
	TestGetBrandList()
	//TestCreateBrandList()
	deleteBrand()
	conn.Close()
}
