package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"shop_api/goods_web/global"
)

func RemoveTopStruct(fields map[string]string) map[string]string {
	rsp := map[string]string{}
	// 遍历传入的字段映射
	for field, err := range fields {
		// 获取字段名中首个"."后的内容，并将其作为新的键，将对应的错误信息作为值存入新的映射中
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

// HandleGrpcErrorToHttp 将 gRPC 错误转换为 HTTP 响应
func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	if err != nil {
		// 判断错误是否为 gRPC 错误
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

// HandleValidatorError 处理验证器错误
func HandleValidatorError(c *gin.Context, err error) {
	// 判断错误是否为验证器错误
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"error": RemoveTopStruct(errs.Translate(global.Trans)),
	})
	return
}

// 以上是一些辅助函数，用于处理错误和字段映射操作。

// 函数 RemoveTopStruct 用于去除字段名中的顶层结构，并返回新的字段映射。
// 它遍历传入的字段映射，通过查找字段名中首个"."的位置，并将"."后的内容作为新的键，
// 将对应的错误信息作为值存入新的映射中，并返回新的映射。

// 函数 HandleGrpcErrorToHttp 用于将 gRPC 错误转换为 HTTP 响应。
// 它首先判断错误是否为 gRPC 错误，如果是，则根据错误码进行相应的处理，
// 并将错误信息和对应的 HTTP 状态码返回给客户端。

// 函数 HandleValidatorError 用于处理验证器错误。
// 它首先判断错误是否为验证器错误，如果是，则将错误信息进行翻译和字段映射处理，
// 并将处理后的错误信息返回给客户端。如果不是验证器错误，则返回原始错误信息。
