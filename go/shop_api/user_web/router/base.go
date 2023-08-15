package router

import (
	"github.com/gin-gonic/gin"
	"shop_api/user_web/api"
)

// InitBaseRouter 初始化基础路由
func InitBaseRouter(Router *gin.RouterGroup) {
	// 创建一个名为BaseRouter的路由组，并将其设置为传入的Router参数的子路由组
	BaseRouter := Router.Group("base")

	{
		// 在BaseRouter中添加一个GET路由，路径为"captcha"，处理函数为api.GetCaptcha
		BaseRouter.GET("captcha", api.GetCaptcha)
	}

}
