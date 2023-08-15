package brands

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"shop_api/goods_web/api"
	"shop_api/goods_web/forms"
	"shop_api/goods_web/global"
	"shop_api/goods_web/proto"
)

func BrandList(ctx *gin.Context) {
	// 获取请求参数的页码和每页数据量
	pn := ctx.DefaultQuery("pn", "0")        // 默认为第一页
	pnInt, _ := strconv.Atoi(pn)             // 将页码转换为整数
	pSize := ctx.DefaultQuery("psize", "10") // 默认每页显示10条数据
	pSizeInt, _ := strconv.Atoi(pSize)       // 将每页数据量转换为整数

	// 调用商品服务的 BrandList 方法获取品牌列表
	rsp, err := global.GoodsSrvClient.BrandList(context.Background(), &proto.BrandFilterRequest{
		Pages:       int32(pnInt),
		PagePerNums: int32(pSizeInt),
	})

	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx) // 处理 gRPC 错误并返回给客户端
		return
	}

	result := make([]interface{}, 0)      // 存储结果的切片
	reMap := make(map[string]interface{}) // 定义返回结果的 map
	reMap["total"] = rsp.Total            // 将品牌总数存入结果

	// 遍历品牌列表，将品牌信息加入结果中
	for _, value := range rsp.Data[pnInt : pnInt*pSizeInt+pSizeInt] {
		reMap := make(map[string]interface{}) // 创建新的品牌信息 map
		reMap["id"] = value.Id                // 将品牌ID加入 map
		reMap["name"] = value.Name            // 将品牌名称加入 map
		reMap["logo"] = value.Logo            // 将品牌logo加入 map

		result = append(result, reMap) // 将品牌信息加入结果切片
	}

	reMap["data"] = result // 将结果切片加入结果 map

	ctx.JSON(http.StatusOK, reMap) // 返回结果给客户端
}

func NewBrand(ctx *gin.Context) {
	brandForm := forms.BrandForm{} // 创建品牌表单对象
	// 绑定请求参数到品牌表单对象
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		api.HandleValidatorError(ctx, err) // 处理验证错误并返回给客户端
		return
	}

	// 调用商品服务的 CreateBrand 方法创建新的品牌
	rsp, err := global.GoodsSrvClient.CreateBrand(context.Background(), &proto.BrandRequest{
		Name: brandForm.Name, // 设置品牌名称
		Logo: brandForm.Logo, // 设置品牌logo
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx) // 处理 gRPC 错误并返回给客户端
		return
	}

	request := make(map[string]interface{}) // 定义返回结果的 map
	request["id"] = rsp.Id                  // 将新品牌的ID加入结果 map
	request["name"] = rsp.Name              // 将新品牌的名称加入结果 map
	request["logo"] = rsp.Logo              // 将新品牌的logo加入结果 map

	ctx.JSON(http.StatusOK, request) // 返回结果给客户端
}

func DeleteBrand(ctx *gin.Context) {
	id := ctx.Param("id")                  // 获取要删除的品牌ID
	i, err := strconv.ParseInt(id, 10, 32) // 将品牌ID转换为整数
	if err != nil {
		ctx.Status(http.StatusNotFound) // 返回404状态码给客户端
		return
	}
	// 调用商品服务的 DeleteBrand 方法删除品牌
	_, err = global.GoodsSrvClient.DeleteBrand(context.Background(), &proto.BrandRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx) // 处理 gRPC 错误并返回给客户端
		return
	}

	ctx.Status(http.StatusOK) // 返回200状态码给客户端
}

func UpdateBrand(ctx *gin.Context) {
	brandForm := forms.BrandForm{} // 创建品牌表单对象
	// 绑定请求参数到品牌表单对象
	if err := ctx.ShouldBindJSON(&brandForm); err != nil {
		api.HandleValidatorError(ctx, err) // 处理验证错误并返回给客户端
		return
	}

	id := ctx.Param("id")                  // 获取要更新的品牌ID
	i, err := strconv.ParseInt(id, 10, 32) // 将品牌ID转换为整数
	if err != nil {
		ctx.Status(http.StatusNotFound) // 返回404状态码给客户端
		return
	}

	// 调用商品服务的 UpdateBrand 方法更新品牌
	_, err = global.GoodsSrvClient.UpdateBrand(context.Background(), &proto.BrandRequest{
		Id:   int32(i),       // 设置要更新的品牌ID
		Name: brandForm.Name, // 设置更新后的品牌名称
		Logo: brandForm.Logo, // 设置更新后的品牌logo
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx) // 处理 gRPC 错误并返回给客户端
		return
	}
	ctx.Status(http.StatusOK) // 返回200状态码给客户端
}

// 获取分类品牌列表
func GetCategoryBrandList(ctx *gin.Context) {
	// 从URL参数中获取id值
	id := ctx.Param("id")
	// 将id转换为int类型
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		// 如果转换失败，则返回404 Not Found错误
		ctx.Status(http.StatusNotFound)
		return
	}

	// 调用Goods服务的GetCategoryBrandList方法获取分类品牌列表数据
	rsp, err := global.GoodsSrvClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: int32(i),
	})
	if err != nil {
		// 如果调用出错，则将gRPC错误转换为HTTP错误并返回
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 构建返回结果数组
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["name"] = value.Name
		reMap["logo"] = value.Logo

		result = append(result, reMap)
	}

	// 返回结果数组
	ctx.JSON(http.StatusOK, result)
}

// 获取分类品牌列表
func CategoryBrandList(ctx *gin.Context) {
	// 调用Goods服务的CategoryBrandList方法获取分类品牌列表数据
	rsp, err := global.GoodsSrvClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{})
	if err != nil {
		// 如果调用出错，则将gRPC错误转换为HTTP错误并返回
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 构建返回结果对象
	reMap := map[string]interface{}{
		"total": rsp.Total,
	}

	// 构建返回结果数组
	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		reMap := make(map[string]interface{})
		reMap["id"] = value.Id
		reMap["category"] = map[string]interface{}{
			"id":   value.Category.Id,
			"name": value.Category.Name,
		}
		reMap["brand"] = map[string]interface{}{
			"id":   value.Brand.Id,
			"name": value.Brand.Name,
			"logo": value.Brand.Logo,
		}

		result = append(result, reMap)
	}

	// 将结果数组添加到结果对象中
	reMap["data"] = result
	// 返回结果对象
	ctx.JSON(http.StatusOK, reMap)
}

// 创建分类品牌
func NewCategoryBrand(ctx *gin.Context) {
	// 从请求中解析出分类品牌表单数据
	categoryBrandForm := forms.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(&categoryBrandForm); err != nil {
		// 如果解析出错，则返回表单验证错误
		api.HandleValidatorError(ctx, err)
		return
	}

	// 调用Goods服务的CreateCategoryBrand方法创建分类品牌
	rsp, err := global.GoodsSrvClient.CreateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		// 如果调用出错，则将gRPC错误转换为HTTP错误并返回
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 构建响应对象
	response := make(map[string]interface{})
	response["id"] = rsp.Id

	// 返回响应对象
	ctx.JSON(http.StatusOK, response)
}

// 更新分类品牌
func UpdateCategoryBrand(ctx *gin.Context) {
	// 从请求中解析出分类品牌表单数据
	categoryBrandForm := forms.CategoryBrandForm{}
	if err := ctx.ShouldBindJSON(&categoryBrandForm); err != nil {
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

	// 调用Goods服务的UpdateCategoryBrand方法更新分类品牌
	_, err = global.GoodsSrvClient.UpdateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		Id:         int32(i),
		CategoryId: int32(categoryBrandForm.CategoryId),
		BrandId:    int32(categoryBrandForm.BrandId),
	})
	if err != nil {
		// 如果调用出错，则将gRPC错误转换为HTTP错误并返回
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.Status(http.StatusOK)
}

// 删除分类品牌
func DeleteCategoryBrand(ctx *gin.Context) {
	// 从URL参数中获取id值
	id := ctx.Param("id")
	// 将id转换为int类型
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		// 如果转换失败，则返回404 Not Found错误
		ctx.Status(http.StatusNotFound)
		return
	}
	// 调用Goods服务的DeleteCategoryBrand方法删除分类品牌
	_, err = global.GoodsSrvClient.DeleteCategoryBrand(context.Background(), &proto.CategoryBrandRequest{Id: int32(i)})
	if err != nil {
		// 如果调用出错，则将gRPC错误转换为HTTP错误并返回
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 返回空对象
	ctx.JSON(http.StatusOK, "")
}
