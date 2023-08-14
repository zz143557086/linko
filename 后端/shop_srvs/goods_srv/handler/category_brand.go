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

func (s *GoodsServer) CategoryBrandList(ctx context.Context, req *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	categoryBrands := []model.GoodsCategoryBrand{}                 // 声明一个空的商品分类品牌列表
	categoryBrandListResponse := proto.CategoryBrandListResponse{} // 声明一个空的商品分类品牌列表响应

	var total int64
	global.DB.Model(&model.GoodsCategoryBrand{}).Count(&total) // 查询商品分类品牌表的总数并保存到total变量中
	categoryBrandListResponse.Total = int32(total)             // 将总数赋值给响应对象的Total字段

	// 根据分页信息查询商品分类品牌列表，同时预加载关联的Category和Brands数据
	global.DB.Preload("Category").Preload("Brands").Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&categoryBrands)

	categoryResponses := []*proto.CategoryBrandResponse{} // 声明一个空的商品分类品牌响应列表
	for _, categoryBrand := range categoryBrands {
		// 创建商品分类品牌响应对象
		categoryResponse := &proto.CategoryBrandResponse{
			Category: &proto.CategoryInfoResponse{
				Id:             categoryBrand.Category.ID,               // 设置分类ID
				Name:           categoryBrand.Category.Name,             // 设置分类名称
				Level:          categoryBrand.Category.Level,            // 设置分类级别
				IsTab:          categoryBrand.Category.IsTab,            // 设置是否在Tab中展示
				ParentCategory: categoryBrand.Category.ParentCategoryID, // 设置父类目ID
			},
			Brand: &proto.BrandInfoResponse{
				Id:   categoryBrand.Brands.ID,   // 设置品牌ID
				Name: categoryBrand.Brands.Name, // 设置品牌名称
				Logo: categoryBrand.Brands.Logo, // 设置品牌Logo
			},
		}
		categoryResponses = append(categoryResponses, categoryResponse) // 将商品分类品牌响应对象添加到列表中
	}

	categoryBrandListResponse.Data = categoryResponses // 将商品分类品牌响应列表赋值给响应对象的Data字段
	return &categoryBrandListResponse, nil             // 返回商品分类品牌列表响应
}

// GetCategoryBrandList 获取商品分类下的品牌列表
func (s *GoodsServer) GetCategoryBrandList(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	brandListResponse := proto.BrandListResponse{}

	// 查询商品分类信息
	var category model.Category
	if result := global.DB.Find(&category, req.Id).First(&category); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	// 查询商品分类下的品牌信息
	var categoryBrands []model.GoodsCategoryBrand
	if result := global.DB.Preload("Brands").Where(&model.GoodsCategoryBrand{CategoryID: req.Id}).Find(&categoryBrands); result.RowsAffected > 0 {
		brandListResponse.Total = int32(result.RowsAffected)
	}

	// 构建品牌信息响应列表
	var brandInfoResponses []*proto.BrandInfoResponse
	for _, categoryBrand := range categoryBrands {
		brandInfoResponses = append(brandInfoResponses, &proto.BrandInfoResponse{
			Id:   categoryBrand.Brands.ID,   // 品牌ID
			Name: categoryBrand.Brands.Name, // 品牌名称
			Logo: categoryBrand.Brands.Logo, // 品牌Logo
		})
	}

	brandListResponse.Data = brandInfoResponses

	return &brandListResponse, nil
}

// CreateCategoryBrand 创建商品分类与品牌的关联关系
func (s *GoodsServer) CreateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	// 查询商品分类信息
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	// 查询品牌信息
	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	// 创建商品分类与品牌的关联关系
	categoryBrand := model.GoodsCategoryBrand{
		CategoryID: req.CategoryId, // 商品分类ID
		BrandsID:   req.BrandId,    // 品牌ID
	}
	global.DB.Save(&categoryBrand)

	return &proto.CategoryBrandResponse{Id: categoryBrand.ID}, nil
}

func (s *GoodsServer) DeleteCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	// 根据请求中的 Id 来删除商品品牌分类
	if result := global.DB.Delete(&model.GoodsCategoryBrand{}, req.Id); result.RowsAffected == 0 {
		// 如果删除操作没有影响任何行（即找不到要删除的品牌分类），返回品牌分类不存在的错误
		return nil, status.Errorf(codes.NotFound, "品牌分类不存在")
	}
	// 返回空的响应和 nil 错误
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategoryBrand(ctx context.Context, req *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	var categoryBrand model.GoodsCategoryBrand

	// 根据请求中的 Id 查询要更新的品牌分类信息
	if result := global.DB.First(&categoryBrand, req.Id); result.RowsAffected == 0 {
		// 如果查询操作没有影响任何行（即找不到要更新的品牌分类），返回品牌分类不存在的错误
		return nil, status.Errorf(codes.InvalidArgument, "品牌分类不存在")
	}

	var category model.Category
	// 根据请求中的 CategoryId 查询对应的商品分类信息
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		// 如果查询操作没有影响任何行（即找不到对应的商品分类），返回商品分类不存在的错误
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	// 根据请求中的 BrandId 查询对应的品牌信息
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		// 如果查询操作没有影响任何行（即找不到对应的品牌），返回品牌不存在的错误
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	// 更新品牌分类的 CategoryID 和 BrandsID 字段
	categoryBrand.CategoryID = req.CategoryId
	categoryBrand.BrandsID = req.BrandId

	// 将更新后的品牌分类信息保存到数据库中
	global.DB.Save(&categoryBrand)

	// 返回空的响应和 nil 错误
	return &emptypb.Empty{}, nil
}
