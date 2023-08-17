package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"shop_api/user_web/api"
	"shop_api/user_web/middlewares"
)

// InitUserRouter 初始化用户路由
func InitUserRouter(Router *gin.RouterGroup) {
	// 创建一个名为UserRouter的路由组，并将其设置为传入的Router参数的子路由组
	UserRouter := Router.Group("/user")

	// 使用zap包记录一条调试级别的日志，表示正在配置用户相关的URL
	zap.S().Debug("配置用户相关的url")

	{
		// 在UserRouter中添加一个GET路由，路径为"list"，处理函数为api.GetUserList
		// 中间件JWTAuth()用于验证JWT，middlewares.IsAdminAuth()用于检查用户权限
		UserRouter.GET("/list", middlewares.JWTAuth() /*, middlewares.IsAdminAuth()*/, api.GetUserList)
		UserRouter.GET("", middlewares.JWTAuth() /*, middlewares.IsAdminAuth()*/, api.GetUserList)
		// 在UserRouter中添加一个POST路由，路径为"login"，处理函数为api.PassWordLogin
		UserRouter.POST("/pwd_login", api.PassWordLogin)

		// 在UserRouter中添加一个POST路由，路径为"register"，处理函数为api.Register
		UserRouter.POST("/register", api.Register)
	}
}
