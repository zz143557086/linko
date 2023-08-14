// ok
package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/emptypb"

	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/model"
	"shop_srvs/goods_srv/proto"
)

// 轮播图列表
func (s *GoodsServer) BannerList(ctx context.Context, req *emptypb.Empty) (*proto.BannerListResponse, error) {
	bannerListResponse := proto.BannerListResponse{}

	// 查询轮播图列表
	var banners []model.Banner
	result := global.DB.Find(&banners) // 执行查询操作

	bannerListResponse.Total = int32(result.RowsAffected)

	// 构建轮播图响应列表
	var bannerResponses []*proto.BannerResponse
	for _, banner := range banners {
		bannerResponses = append(bannerResponses, &proto.BannerResponse{
			Id:    banner.ID,
			Image: banner.Image,
			Index: banner.Index,
			Url:   banner.Url,
		})
	}

	bannerListResponse.Data = bannerResponses

	return &bannerListResponse, nil
}

// 创建轮播图
func (s *GoodsServer) CreateBanner(ctx context.Context, req *proto.BannerRequest) (*proto.BannerResponse, error) {
	banner := model.Banner{}

	// 设置轮播图属性
	banner.Image = req.Image
	banner.Index = req.Index
	banner.Url = req.Url

	// 保存轮播图到数据库
	global.DB.Save(&banner)

	return &proto.BannerResponse{Id: banner.ID}, nil
}

// 删除轮播图
func (s *GoodsServer) DeleteBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	// 根据ID查询轮播图
	if result := global.DB.Delete(&model.Banner{}, req.Id); result.RowsAffected == 0 {
		// 如果轮播图不存在，则返回错误
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}
	return &emptypb.Empty{}, nil
}

// 更新轮播图
func (s *GoodsServer) UpdateBanner(ctx context.Context, req *proto.BannerRequest) (*emptypb.Empty, error) {
	var banner model.Banner

	// 根据ID查询轮播图
	if result := global.DB.First(&banner, req.Id); result.RowsAffected == 0 {
		// 如果轮播图不存在，则返回错误
		return nil, status.Errorf(codes.NotFound, "轮播图不存在")
	}

	// 更新轮播图属性
	if req.Url != "" {
		banner.Url = req.Url
	}
	if req.Image != "" {
		banner.Image = req.Image
	}
	if req.Index != 0 {
		banner.Index = req.Index
	}

	// 保存更新后的轮播图到数据库
	global.DB.Save(&banner)

	return &emptypb.Empty{}, nil
}
