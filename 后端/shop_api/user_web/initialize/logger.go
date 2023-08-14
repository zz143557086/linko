package initialize

import "go.uber.org/zap"

func InitLogger() {
	// 创建一个开发环境下的日志记录器
	logger, _ := zap.NewDevelopment()

	// 将全局的日志记录器替换为新创建的日志记录器
	zap.ReplaceGlobals(logger)

	// 使用日志记录器输出调试日志信息
	zap.S().Debugf("日志加载成功")
}
