package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"shop_api/user_web/models"
)

// IsAdminAuth 是一个 Gin 中间件，用于检查当前用户是否具有管理员权限
func IsAdminAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从上下文中获取声明信息
		claims, _ := ctx.Get("claims")
		currentUser := claims.(*models.CustomClaims)

		// 检查当前用户的权限ID是否为2，若不是，则表示没有管理员权限
		if currentUser.AuthorityId != 2 {
			ctx.JSON(http.StatusForbidden, gin.H{
				"msg": "无权限",
			})
			ctx.Abort()
			return
		}

		// 继续处理后续的请求
		ctx.Next()
	}
}
