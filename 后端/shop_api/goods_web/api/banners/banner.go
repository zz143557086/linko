package banners

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/empty"

	"shop_api/goods_web/api"
	"shop_api/goods_web/forms"
	"shop_api/goods_web/global"
	"shop_api/goods_web/proto"
)

// List 获取轮播图列表
func List(ctx *gin.Context) {
	// 调用Goods服务的BannerList方法获取轮播图列表数据
	rsp, err := global.GoodsSrvClient.BannerList(context.Background(), &empty.Empty{})
	if err != nil {
		// 如果调用出错，则将gRPC错误转换为HTTP错误并返回
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 构建返回结果数组
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		// 创建一个map来存储单个轮播图的信息
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id       // 轮播图ID
		reMap["index"] = value.Index // 轮播图索引
		reMap["image"] = value.Image // 轮播图图片
		reMap["url"] = value.Url     // 轮播图链接

		// 添加到结果数组中
		result = append(result, reMap)
	}

	// 返回结果数组
	ctx.JSON(http.StatusOK, result)
}

// New 创建轮播图
func New(ctx *gin.Context) {
	// 从请求中解析出轮播图表单数据
	bannerForm := forms.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		// 如果解析出错，则返回表单验证错误
		api.HandleValidatorError(ctx, err)
		return
	}

	// 调用Goods服务的CreateBanner方法创建轮播图
	rsp, err := global.GoodsSrvClient.CreateBanner(context.Background(), &proto.BannerRequest{
		Index: int32(bannerForm.Index), // 轮播图索引
		Url:   bannerForm.Url,          // 轮播图链接
		Image: bannerForm.Image,        // 轮播图图片
	})
	if err != nil {
		// 如果调用出错，则将gRPC错误转换为HTTP错误并返回
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 构建响应对象
	response := make(map[string]interface{})
	response["id"] = rsp.Id       // 创建的轮播图ID
	response["index"] = rsp.Index // 创建的轮播图索引
	response["url"] = rsp.Url     // 创建的轮播图链接
	response["image"] = rsp.Image // 创建的轮播图图片

	// 返回响应对象
	ctx.JSON(http.StatusOK, response)
}

// Update 更新轮播图
func Update(ctx *gin.Context) {
	// 从请求中解析出轮播图表单数据
	bannerForm := forms.BannerForm{}
	if err := ctx.ShouldBindJSON(&bannerForm); err != nil {
		// 如果解析出错，则返回表单验证错误
		api.HandleValidatorError(ctx, err)
		return
	}

	// 从URL参数中获取id值
	id := ctx.Param("id")
	// 将id转换为int类型
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		// 如果转换失败，则返回404 Not Found错误
		ctx.Status(http.StatusNotFound)
		return
	}

	// 调用Goods服务的UpdateBanner方法更新轮播图
	_, err = global.GoodsSrvClient.UpdateBanner(context.Background(), &proto.BannerRequest{
		Id:    int32(i),                // 轮播图ID
		Index: int32(bannerForm.Index), // 轮播图索引
		Url:   bannerForm.Url,          // 轮播图链接
	})
	if err != nil {
		// 如果调用出错，则将gRPC错误转换为HTTP错误并返回
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 返回成功状态
	ctx.Status(http.StatusOK)
}

// Delete 删除轮播图
func Delete(ctx *gin.Context) {
	// 从URL参数中获取id值
	id := ctx.Param("id")
	// 将id转换为int类型
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		// 如果转换失败，则返回404 Not Found错误
		ctx.Status(http.StatusNotFound)
		return
	}

	// 调用Goods服务的DeleteBanner方法删除轮播图
	_, err = global.GoodsSrvClient.DeleteBanner(context.Background(), &proto.BannerRequest{Id: int32(i)})
	if err != nil {
		// 如果调用出错，则将gRPC错误转换为HTTP错误并返回
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 返回空对象
	ctx.JSON(http.StatusOK, "")
}
