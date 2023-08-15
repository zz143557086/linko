package goods

import (
	"context"
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"shop_api/goods_web/forms"
	"shop_api/goods_web/proto"
	"strconv"
	"strings"

	"shop_api/goods_web/global"
)

func removeTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	// 遍历错误字段名和错误信息的映射
	for field, err := range fields {
		// 移除顶层结构，只保留字段名
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	// 将 gRPC 错误转换成对应的 HTTP 响应状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg:": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": e.Code(),
				})
			}
			return
		}
	}
}

func HandleValidatorError(c *gin.Context, err error) {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 如果错误不是验证器错误，则返回普通错误信息
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
		return
	}

	// 构造错误信息的映射并返回
	c.JSON(http.StatusBadRequest, gin.H{
		"error": removeTopStruct(errs.Translate(global.Trans)),
	})
	return
}

func List(ctx *gin.Context) {
	fmt.Println("商品列表")
	// 定义商品筛选请求对象
	request := &proto.GoodsFilterRequest{}

	// 获取查询参数 pmin，指定商品最低价格
	priceMin := ctx.DefaultQuery("pmin", "0")
	priceMinInt, _ := strconv.Atoi(priceMin)
	request.PriceMin = int32(priceMinInt)

	// 获取查询参数 pmax，指定商品最高价格
	priceMax := ctx.DefaultQuery("pmax", "0")
	priceMaxInt, _ := strconv.Atoi(priceMax)
	request.PriceMax = int32(priceMaxInt)

	// 获取查询参数 ih，是否热销商品
	isHot := ctx.DefaultQuery("ih", "0")
	if isHot == "1" {
		request.IsHot = true
	}

	// 获取查询参数 in，是否新品商品
	isNew := ctx.DefaultQuery("in", "0")
	if isNew == "1" {
		request.IsNew = true
	}

	// 获取查询参数 it，是否标签商品
	isTab := ctx.DefaultQuery("it", "0")
	if isTab == "1" {
		request.IsTab = true
	}

	// 获取查询参数 c，指定商品分类ID
	categoryId := ctx.DefaultQuery("c", "0")
	categoryIdInt, _ := strconv.Atoi(categoryId)
	request.TopCategory = int32(categoryIdInt)

	// 获取查询参数 p，指定当前页数
	pages := ctx.DefaultQuery("p", "0")
	pagesInt, _ := strconv.Atoi(pages)
	request.Pages = int32(pagesInt)

	// 获取查询参数 pnum，指定每页展示的商品数量
	perNums := ctx.DefaultQuery("pnum", "0")
	perNumsInt, _ := strconv.Atoi(perNums)
	request.PagePerNums = int32(perNumsInt)

	// 获取查询参数 q，指定关键字搜索
	keywords := ctx.DefaultQuery("q", "")
	request.KeyWords = keywords

	// 获取查询参数 b，指定品牌ID
	brandId := ctx.DefaultQuery("b", "0")
	brandIdInt, _ := strconv.Atoi(brandId)
	request.Brand = int32(brandIdInt)

	// 商品服务客户端请求商品列表
	e, b := sentinel.Entry("goods-list", sentinel.WithTrafficType(base.Inbound))
	if b != nil {
		ctx.JSON(http.StatusTooManyRequests, gin.H{
			"msg": "请求过于频繁，请稍后重试",
		})
		return
	}
	r, err := global.GoodsSrvClient.GoodsList(context.WithValue(context.Background(), "ginContext", ctx), request)
	if err != nil {
		zap.S().Errorw("[List] 查询 【商品列表】失败")
		HandleGrpcErrorToHttp(err, ctx)
		return
	}
	e.Exit()

	// 构建响应数据
	reMap := map[string]interface{}{
		"total": r.Total,
	}

	goodsList := make([]interface{}, 0)
	for _, value := range r.Data {
		// 将商品对象转换为可序列化的 map 结构
		goodsList = append(goodsList, map[string]interface{}{
			"id":          value.Id,
			"name":        value.Name,
			"goods_brief": value.GoodsBrief,
			"desc":        value.GoodsDesc,
			"ship_free":   value.ShipFree,
			"images":      value.Images,
			"desc_images": value.DescImages,
			"front_image": value.GoodsFrontImage,
			"shop_price":  value.ShopPrice,
			"category": map[string]interface{}{
				"id":   value.Category.Id,
				"name": value.Category.Name,
			},
			"brand": map[string]interface{}{
				"id":   value.Brand.Id,
				"name": value.Brand.Name,
				"logo": value.Brand.Logo,
			},
			"is_hot":  value.IsHot,
			"is_new":  value.IsNew,
			"on_sale": value.OnSale,
		})
	}
	reMap["data"] = goodsList

	ctx.JSON(http.StatusOK, reMap)
}

// 该方法实现了商品列表的查询功能。
// 首先根据查询参数构造了商品筛选请求对象，包括价格、热销、新品、分类、关键字、品牌等信息。
// 然后调用商品服务的 GoodsList 方法查询商品列表，并进行错误处理。
// 最后将查询结果构造为统一的响应格式返回给客户端。

// New 用于创建商品
func New(ctx *gin.Context) {
	// 解析请求中的商品表单数据
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// 调用商品服务的 CreateGoods 方法创建商品
	goodsClient := global.GoodsSrvClient
	rsp, err := goodsClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Name:            goodsForm.Name,        // 商品名称
		GoodsSn:         goodsForm.GoodsSn,     // 商品编号
		Stocks:          goodsForm.Stocks,      // 库存数量
		MarketPrice:     goodsForm.MarketPrice, // 市场价格
		ShopPrice:       goodsForm.ShopPrice,   // 商店价格
		GoodsBrief:      goodsForm.GoodsBrief,  // 商品简介
		ShipFree:        *goodsForm.ShipFree,   // 是否包邮
		Images:          goodsForm.Images,      // 商品图片列表
		DescImages:      goodsForm.DescImages,  // 商品描述图片列表
		GoodsFrontImage: goodsForm.FrontImage,  // 商品封面图片
		CategoryId:      goodsForm.CategoryId,  // 商品分类ID
		BrandId:         goodsForm.Brand,       // 商品品牌ID
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 返回创建成功的商品信息
	ctx.JSON(http.StatusOK, rsp)
}

// Detail 获取商品详情
func Detail(ctx *gin.Context) {
	// 从 URL 参数中获取商品 ID
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// 调用商品服务的 GetGoodsDetail 方法获取商品详情
	r, err := global.GoodsSrvClient.GetGoodsDetail(context.WithValue(context.Background(), "ginContext", ctx), &proto.GoodInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 构建商品详情的响应数据
	rsp := map[string]interface{}{
		"id":          r.Id,              // 商品ID
		"name":        r.Name,            // 商品名称
		"goods_brief": r.GoodsBrief,      // 商品简介
		"desc":        r.GoodsDesc,       // 商品描述
		"ship_free":   r.ShipFree,        // 是否包邮
		"images":      r.Images,          // 商品图片列表
		"desc_images": r.DescImages,      // 商品描述图片列表
		"front_image": r.GoodsFrontImage, // 商品封面图片
		"shop_price":  r.ShopPrice,       // 商店价格
		"category": map[string]interface{}{ // 商品分类信息
			"id":   r.Category.Id,   // 分类ID
			"name": r.Category.Name, // 分类名称
		},
		"brand": map[string]interface{}{ // 商品品牌信息
			"id":   r.Brand.Id,   // 品牌ID
			"name": r.Brand.Name, // 品牌名称
			"logo": r.Brand.Logo, // 品牌Logo
		},
		"is_hot":  r.IsHot,  // 是否热销商品
		"is_new":  r.IsNew,  // 是否新品上市
		"on_sale": r.OnSale, // 是否在售
	}

	// 返回商品详情
	ctx.JSON(http.StatusOK, rsp)
}

// Delete 删除商品
func Delete(ctx *gin.Context) {
	// 从 URL 参数中获取商品 ID
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// 调用商品服务的 DeleteGoods 方法删除商品
	_, err = global.GoodsSrvClient.DeleteGoods(context.Background(), &proto.DeleteGoodsInfo{Id: int32(i)})
	if err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 返回成功状态码
	ctx.Status(http.StatusOK)
	return
}

// Stocks 获取商品库存
func Stocks(ctx *gin.Context) {
	// 从 URL 参数中获取商品 ID
	id := ctx.Param("id")
	_, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// TODO: 根据商品 ID 获取商品库存

	return
}

// UpdateStatus 更新商品状态
func UpdateStatus(ctx *gin.Context) {
	// 解析请求中的商品状态表单数据
	goodsStatusForm := forms.GoodsStatusForm{}
	if err := ctx.ShouldBindJSON(&goodsStatusForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// 从 URL 参数中获取商品 ID
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:     int32(i),
		IsHot:  *goodsStatusForm.IsHot,  // 是否热销商品
		IsNew:  *goodsStatusForm.IsNew,  // 是否新品上市
		OnSale: *goodsStatusForm.OnSale, // 是否在售
	}); err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

// Update 更新商品信息
func Update(ctx *gin.Context) {
	// 解析请求中的商品表单数据
	goodsForm := forms.GoodsForm{}
	if err := ctx.ShouldBindJSON(&goodsForm); err != nil {
		HandleValidatorError(ctx, err)
		return
	}

	// 从 URL 参数中获取商品 ID
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if _, err = global.GoodsSrvClient.UpdateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              int32(i),
		Name:            goodsForm.Name,        // 商品名称
		GoodsSn:         goodsForm.GoodsSn,     // 商品编号
		Stocks:          goodsForm.Stocks,      // 库存数量
		MarketPrice:     goodsForm.MarketPrice, // 市场价格
		ShopPrice:       goodsForm.ShopPrice,   // 商店价格
		GoodsBrief:      goodsForm.GoodsBrief,  // 商品简介
		ShipFree:        *goodsForm.ShipFree,   // 是否包邮
		Images:          goodsForm.Images,      // 商品图片列表
		DescImages:      goodsForm.DescImages,  // 商品描述图片列表
		GoodsFrontImage: goodsForm.FrontImage,  // 商品封面图片
		CategoryId:      goodsForm.CategoryId,  // 商品分类ID
		BrandId:         goodsForm.Brand,       // 商品品牌ID
	}); err != nil {
		HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "更新成功",
	})
}
