package main

import (
	"fmt"
	"github.com/0xweb-3/CoinNest/eth_srv/global"
	"github.com/0xweb-3/CoinNest/eth_srv/handler"
	"github.com/0xweb-3/CoinNest/eth_srv/initialize"
	"github.com/0xweb-3/CoinNest/eth_srv/service"
	"github.com/0xweb-3/CoinNest/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. 初始化日志
	initialize.InitLogger()

	// 2. 初始化配置信息
	initialize.InitConfig()

	// 3. 初始化数据库
	initialize.InitDB()

	IP := global.ServerConfig.Host
	Port := global.ServerConfig.Port

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", IP, Port))
	if err != nil {
		zap.S().Fatalf("failed to listen: %s", err.Error())
	}
	s := grpc.NewServer()
	//  注册服务
	ethRepo := handler.NewUserRepo(global.DB)
	srv := service.NewEthServer(ethRepo)
	proto.RegisterUserServer(s, srv)

	reflection.Register(s)
	// 启动服务
	zap.S().Debugf("service listening at: %v", lis.Addr())

	go func() {
		if err := s.Serve(lis); err != nil {
			zap.S().Fatalf("failed to serve: %v", err)
		}
	}()

	// 接受服务退出信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.S().Info("Shutting down GRPC server...")
	s.GracefulStop()
}
