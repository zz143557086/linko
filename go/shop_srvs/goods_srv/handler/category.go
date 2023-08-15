package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"shop_srvs/goods_srv/model"

	"google.golang.org/protobuf/types/known/emptypb"
	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/proto"
)

// 商品分类
func (s *GoodsServer) GetAllCategorysList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	/*
		获取所有商品分类列表
		返回数据格式:
		[
			{
				"id":xxx,
				"name":"",
				"level":1,
				"is_tab":false,
				"parent":13xxx,
				"sub_category":[
					{
						"id":xxx,
						"name":"",
						"level":1,
						"is_tab":false,
						"sub_category":[]
					}
				]
			}
		]
	*/
	// 查询所有一级分类
	var categorys []model.Category
	global.DB.Where(&model.Category{Level: 1}).Preload("SubCategory.SubCategory").Find(&categorys)
	b, _ := json.Marshal(&categorys)
	return &proto.CategoryListResponse{JsonData: string(b)}, nil
}

// 获取子分类
func (s *GoodsServer) GetSubCategory(ctx context.Context, req *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	categoryListResponse := proto.SubCategoryListResponse{}

	// 查询父分类
	var category model.Category
	if result := global.DB.Where("id=?", req.Id).First(&category); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	categoryListResponse.Info = &proto.CategoryInfoResponse{
		Id:             category.ID,               // 父分类的ID
		Name:           category.Name,             // 父分类的名称
		Level:          category.Level,            // 父分类的层级
		IsTab:          category.IsTab,            // 父分类是否为标签分类
		ParentCategory: category.ParentCategoryID, // 父分类的父分类ID（如果有）
	}

	// 查询子分类
	var subCategorys []model.Category
	var subCategoryResponse []*proto.CategoryInfoResponse
	global.DB.Where(&model.Category{ParentCategoryID: req.Id}).Find(&subCategorys)
	for _, subCategory := range subCategorys {
		subCategoryResponse = append(subCategoryResponse, &proto.CategoryInfoResponse{
			Id:             subCategory.ID,               // 子分类的ID
			Name:           subCategory.Name,             // 子分类的名称
			Level:          subCategory.Level,            // 子分类的层级
			IsTab:          subCategory.IsTab,            // 子分类是否为标签分类
			ParentCategory: subCategory.ParentCategoryID, // 子分类的父分类ID
		})
	}

	categoryListResponse.SubCategorys = subCategoryResponse
	return &categoryListResponse, nil
}
func (s *GoodsServer) CreateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	category := model.Category{}
	cMap := map[string]interface{}{}
	cMap["name"] = req.Name    // 分类名称
	cMap["level"] = req.Level  // 分类级别
	cMap["is_tab"] = req.IsTab // 是否在Tab中展示
	if req.Level != 1 {
		// 查询父类目是否存在
		cMap["parent_category_id"] = req.ParentCategory // 父类目ID
	}
	tx := global.DB.Model(&model.Category{}).Create(cMap)
	fmt.Println(tx)
	return &proto.CategoryInfoResponse{Id: int32(category.ID)}, nil
}

func (s *GoodsServer) DeleteCategory(ctx context.Context, req *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Category{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}
	return &emptypb.Empty{}, nil
}

func (s *GoodsServer) UpdateCategory(ctx context.Context, req *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	var category model.Category

	if result := global.DB.First(&category, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品分类不存在")
	}

	if req.Name != "" {
		category.Name = req.Name // 更新分类名称
	}
	if req.ParentCategory != 0 {
		category.ParentCategoryID = req.ParentCategory // 更新父类目ID
	}
	if req.Level != 0 {
		category.Level = req.Level // 更新分类级别
	}
	if req.IsTab {
		category.IsTab = req.IsTab // 更新是否在Tab中展示
	}

	global.DB.Save(&category)

	return &emptypb.Empty{}, nil
}
