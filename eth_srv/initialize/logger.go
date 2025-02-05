package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
	// 日志级别：debug info warn error fetal
}
