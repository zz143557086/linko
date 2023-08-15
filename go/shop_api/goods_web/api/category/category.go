package category

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	empty "github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"

	"shop_api/goods_web/api"
	"shop_api/goods_web/forms"
	"shop_api/goods_web/global"
	"shop_api/goods_web/proto"
)

// List 获取所有分类列表
func List(ctx *gin.Context) {
	// 调用商品服务的 GetAllCategorysList 方法获取分类列表
	r, err := global.GoodsSrvClient.GetAllCategorysList(context.Background(), &empty.Empty{})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	data := make([]interface{}, 0)
	err = json.Unmarshal([]byte(r.JsonData), &data)
	if err != nil {
		zap.S().Errorw("[List] 查询 【分类列表】失败： ", err.Error())
	}

	// 返回分类列表数据
	ctx.JSON(http.StatusOK, data)
}

// Detail 获取分类详情
func Detail(ctx *gin.Context) {
	// 从 URL 参数中获取分类 ID
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	reMap := make(map[string]interface{})
	subCategorys := make([]interface{}, 0)
	if r, err := global.GoodsSrvClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id: int32(i),
	}); err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	} else {
		// 对子分类列表进行处理
		for _, value := range r.SubCategorys {
			subCategorys = append(subCategorys, map[string]interface{}{
				"id":              value.Id,             // 子分类ID
				"name":            value.Name,           // 子分类名称
				"level":           value.Level,          // 子分类层级
				"parent_category": value.ParentCategory, // 父分类ID
				"is_tab":          value.IsTab,          // 是否在导航栏显示
			})
		}

		// 设置分类详情数据
		reMap["id"] = r.Info.Id                          // 分类ID
		reMap["name"] = r.Info.Name                      // 分类名称
		reMap["level"] = r.Info.Level                    // 分类层级
		reMap["parent_category"] = r.Info.ParentCategory // 父分类ID
		reMap["is_tab"] = r.Info.IsTab                   // 是否在导航栏显示
		reMap["sub_categorys"] = subCategorys            // 子分类列表

		ctx.JSON(http.StatusOK, reMap)
	}
	return
}

// New 创建分类
func New(ctx *gin.Context) {
	categoryForm := forms.CategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 调用商品服务的 CreateCategory 方法创建分类
	rsp, err := global.GoodsSrvClient.CreateCategory(context.Background(), &proto.CategoryInfoRequest{
		Name:           categoryForm.Name,           // 分类名称
		IsTab:          *categoryForm.IsTab,         // 是否在导航栏显示
		Level:          categoryForm.Level,          // 分类层级
		ParentCategory: categoryForm.ParentCategory, // 父分类ID
	})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	request := make(map[string]interface{})
	request["id"] = rsp.Id                 // 分类ID
	request["name"] = rsp.Name             // 分类名称
	request["parent"] = rsp.ParentCategory // 父分类ID
	request["level"] = rsp.Level           // 分类层级
	request["is_tab"] = rsp.IsTab          // 是否在导航栏显示

	// 返回创建的分类信息
	ctx.JSON(http.StatusOK, request)
}

// Delete 删除分类
func Delete(ctx *gin.Context) {
	// 从 URL 参数中获取分类 ID
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	// 1. 先查询出该分类下的所有子分类
	// 2. 将所有的分类全部逻辑删除
	// 3. 将该分类下的所有的商品逻辑删除
	_, err = global.GoodsSrvClient.DeleteCategory(context.Background(), &proto.DeleteCategoryRequest{Id: int32(i)})
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 返回成功状态码
	ctx.Status(http.StatusOK)
}

// Update 更新分类信息
func Update(ctx *gin.Context) {
	categoryForm := forms.UpdateCategoryForm{}
	if err := ctx.ShouldBindJSON(&categoryForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	// 从 URL 参数中获取分类 ID
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	request := &proto.CategoryInfoRequest{
		Id:   int32(i),          // 分类ID
		Name: categoryForm.Name, // 更新后的分类名称
	}
	if categoryForm.IsTab != nil {
		request.IsTab = *categoryForm.IsTab // 更新后的是否在导航栏显示
	}

	// 调用商品服务的 UpdateCategory 方法更新分类信息
	_, err = global.GoodsSrvClient.UpdateCategory(context.Background(), request)
	if err != nil {
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	// 返回成功状态码
	ctx.Status(http.StatusOK)
}
