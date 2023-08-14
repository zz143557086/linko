package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Cors 是一个 Gin 中间件，用于处理跨域请求
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// 设置跨域请求的响应头
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			// 对于预检请求，直接返回成功的响应，并终止请求处理
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
