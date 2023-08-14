package initialize

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop_api/user_web/middlewares"
	"shop_api/user_web/router"
)

func Routers() *gin.Engine {
	// 创建一个默认的gin引擎作为路由器
	Router := gin.Default()

	// 添加跨域中间件
	Router.Use(middlewares.Cors())

	// 创建一个路由组，前缀为"/u/v1"
	ApiGroup := Router.Group("/u/v1")

	// 初始化用户相关的路由
	router.InitUserRouter(ApiGroup)
	zap.S().Debug("用户路由初始化成功")

	// 初始化基本信息相关的路由
	router.InitBaseRouter(ApiGroup)
	zap.S().Debug("基本信息初始化成功")

	return Router
}
