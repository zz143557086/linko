package handler

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shop_srvs/goods_srv/global"
	"shop_srvs/goods_srv/model"
	"shop_srvs/goods_srv/proto"
)

type GoodsServer struct {
	proto.UnimplementedGoodsServer
}

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	// 将商品模型转换为响应模型
	// 将模型对象中的字段拷贝到响应对象中
	return proto.GoodsInfoResponse{
		Id:              goods.ID,              // 商品ID
		CategoryId:      goods.CategoryID,      // 商品分类ID
		Name:            goods.Name,            // 商品名称
		GoodsSn:         goods.GoodsSn,         // 商品货号
		ClickNum:        goods.ClickNum,        // 商品点击量
		SoldNum:         goods.SoldNum,         // 商品销售量
		FavNum:          goods.FavNum,          // 商品收藏量
		MarketPrice:     goods.MarketPrice,     // 商品市场价
		ShopPrice:       goods.ShopPrice,       // 商品商城价
		GoodsBrief:      goods.GoodsBrief,      // 商品简介
		ShipFree:        goods.ShipFree,        // 是否免费配送
		GoodsFrontImage: goods.GoodsFrontImage, // 商品封面图
		IsNew:           goods.IsNew,           // 是否新品
		IsHot:           goods.IsHot,           // 是否热销
		OnSale:          goods.OnSale,          // 是否在售
		DescImages:      goods.DescImages,      // 商品描述图片
		Images:          goods.Images,          // 商品图片
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,   // 商品分类ID
			Name: goods.Category.Name, // 商品分类名称
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,   // 品牌ID
			Name: goods.Brands.Name, // 品牌名称
			Logo: goods.Brands.Logo, // 品牌Logo
		},
	}
}

func (s *GoodsServer) GoodsList(ctx context.Context, req *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	// 商品列表接口实现
	// 根据传入的过滤条件搜索和筛选商品
	goodsListResponse := &proto.GoodsListResponse{}

	// 使用本地数据库进行查询
	localDB := global.DB.Model(model.Goods{})

	if req.KeyWords != "" {
		// 如果关键词不为空，则将关键词应用于商品名称和商品简介字段的多字段查询
		localDB.Where("name LIKE ? ", "%"+req.KeyWords+"%")
	}

	if req.IsHot {
		// 如果需要筛选热销商品，则在本地数据库中添加过滤条件
		localDB = localDB.Where(model.Goods{IsHot: true})
	}

	if req.IsNew {
		// 如果需要筛选新品，则在布尔查询中添加过滤条件
		localDB = localDB.Where(model.Goods{IsNew: true})
	}

	if req.PriceMin > 0 {
		// 如果设置了最低价格，则在布尔查询中添加价格范围过滤条件
		localDB = localDB.Where("shop_price >= ? ", req.PriceMin)
	}

	if req.PriceMax > 0 {
		// 如果设置了最高价格，则在布尔查询中添加价格范围过滤条件
		localDB = localDB.Where("shop_price <= ? ", req.PriceMax)
	}

	if req.Brand > 0 {
		// 如果设置了品牌ID，则在布尔查询中添加品牌ID过滤条件
		localDB = localDB.Where("bind_id = ? ", req.Brand)
	}

	// 通过分类查询商品
	var subQuery string
	//categoryIds := make([]interface{}, 0)

	if req.TopCategory > 0 {
		// 如果设置了顶级分类ID，则根据分类层级构造子查询语句
		var category model.Category

		// 根据顶级分类ID查询分类信息
		if result := global.DB.First(&category, req.TopCategory); result.RowsAffected == 0 {
			return nil, status.Errorf(codes.NotFound, "商品分类不存在")
		}

		// 根据分类层级生成子查询语句
		switch category.Level {
		case 1:
			// 如果顶级分类层级为1，说明需要查询二级分类下的商品
			subQuery = fmt.Sprintf("select id from category where parent_category_id in (select id from category WHERE parent_category_id=%d)", req.TopCategory)
		case 2:
			// 如果顶级分类层级为2，说明需要查询三级分类下的商品
			subQuery = fmt.Sprintf("select id from category WHERE parent_category_id=%d", req.TopCategory)
		case 3:
			// 如果顶级分类层级为3，说明需要查询当前分类下的商品
			subQuery = fmt.Sprintf("select id from category WHERE id=%d", req.TopCategory)
		}
		localDB = localDB.Where(fmt.Sprintf("category_id in (%s)", subQuery))
	}

	// 分页处理
	var count int64
	var goods []model.Goods
	localDB.Count(&count)
	goodsListResponse.Total = int32(count)
	result := localDB.Preload("Category").Preload("Brands").
		Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&goods)
	if result.Error != nil {
		return nil, result.Error
	}

	// 将查询到的商品信息转换为响应模型，并添加到响应对象中

	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	// 返回商品列表响应
	return goodsListResponse, nil
}

// BatchGetGoods 批量查询商品信息
func (s *GoodsServer) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods

	// 根据商品ID列表查询商品信息
	result := global.DB.Where(req.Id).Find(&goods)
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}

	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
}

// GetGoodsDetail 获取商品详情
func (s *GoodsServer) GetGoodsDetail(ctx context.Context, req *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	var goods model.Goods

	// 根据商品ID查询商品详情，并预加载关联的分类和品牌信息
	if result := global.DB.Preload("Category").Preload("Brands").First(&goods, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}

	goodsInfoResponse := ModelToResponse(goods)
	return &goodsInfoResponse, nil
}

// CreateGoods 创建商品
func (s *GoodsServer) CreateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	var category model.Category
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	// 创建商品对象
	goods := model.Goods{
		Brands:          brand,
		BrandsID:        brand.ID,
		Category:        category,
		CategoryID:      category.ID,
		Name:            req.Name,            // 商品名称
		GoodsSn:         req.GoodsSn,         // 商品编号
		MarketPrice:     req.MarketPrice,     // 市场价格
		ShopPrice:       req.ShopPrice,       // 商城价格
		GoodsBrief:      req.GoodsBrief,      // 商品简介
		ShipFree:        req.ShipFree,        // 是否包邮
		Images:          req.Images,          // 商品图片列表
		DescImages:      req.DescImages,      // 商品详情图片列表
		GoodsFrontImage: req.GoodsFrontImage, // 商品封面图片
		IsNew:           req.IsNew,           // 是否新品
		IsHot:           req.IsHot,           // 是否热卖
		OnSale:          req.OnSale,          // 是否上架
	}

	// 在全局的数据库事务中执行保存操作
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()

	return &proto.GoodsInfoResponse{
		Id: goods.ID, // 商品ID
	}, nil
}

// DeleteGoods 删除商品
func (s *GoodsServer) DeleteGoods(ctx context.Context, req *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	// 根据商品ID删除商品，并使用模型的基本字段来设置商品ID
	if result := global.DB.Delete(&model.Goods{BaseModel: model.BaseModel{ID: req.Id}}, req.Id); result.Error != nil {
		// 如果删除操作失败，则返回商品不存在的错误
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	return &emptypb.Empty{}, nil
}

// UpdateGoods 更新商品信息
func (s *GoodsServer) UpdateGoods(ctx context.Context, req *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	var goods model.Goods

	// 根据商品ID查询商品信息
	if result := global.DB.First(&goods, req.Id); result.RowsAffected == 0 {
		// 如果商品不存在，则返回商品不存在的错误
		return nil, status.Errorf(codes.NotFound, "商品不存在")
	}
	if req.CategoryId == 0 || req.BrandId == 0 {
		// 在全局的数据库事务中执行保存操作
		tx := global.DB.Begin()
		result := tx.Model(&model.Goods{}).Where("id = ?", goods.ID).
			Updates(model.Goods{IsNew: req.IsNew, IsHot: req.IsHot, OnSale: req.OnSale}) // 更新部分字段
		if result.Error != nil {
			tx.Rollback()
			return nil, result.Error
		}
		tx.Commit()
		return &emptypb.Empty{}, nil

	}

	var category model.Category
	// 根据商品分类ID查询商品分类信息
	if result := global.DB.First(&category, req.CategoryId); result.RowsAffected == 0 {
		// 如果商品分类不存在，则返回商品分类不存在的错误
		return nil, status.Errorf(codes.InvalidArgument, "商品分类不存在")
	}

	var brand model.Brands
	// 根据品牌ID查询品牌信息
	if result := global.DB.First(&brand, req.BrandId); result.RowsAffected == 0 {
		// 如果品牌不存在，则返回品牌不存在的错误
		return nil, status.Errorf(codes.InvalidArgument, "品牌不存在")
	}

	// 更新商品信息
	goods.Brands = brand
	goods.BrandsID = brand.ID
	goods.Category = category
	goods.CategoryID = category.ID
	goods.Name = req.Name                       // 商品名称
	goods.GoodsSn = req.GoodsSn                 // 商品编号
	goods.MarketPrice = req.MarketPrice         // 市场价格
	goods.ShopPrice = req.ShopPrice             // 商城价格
	goods.GoodsBrief = req.GoodsBrief           // 商品简介
	goods.ShipFree = req.ShipFree               // 是否包邮
	goods.Images = req.Images                   // 商品图片列表
	goods.DescImages = req.DescImages           // 商品详情图片列表
	goods.GoodsFrontImage = req.GoodsFrontImage // 商品封面图片
	goods.IsNew = req.IsNew                     // 是否新品
	goods.IsHot = req.IsHot                     // 是否热卖
	goods.OnSale = req.OnSale                   // 是否上架

	// 在全局的数据库事务中执行保存操作
	tx := global.DB.Begin()
	result := tx.Save(&goods)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
