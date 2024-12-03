package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/user_srv/global"
	"mxshop_srvs/user_srv/handler"
	"mxshop_srvs/user_srv/initialize"
	"mxshop_srvs/user_srv/proto"
	"mxshop_srvs/user_srv/utils"
)

func main() {
	// 1、初始化zap
	initialize.InitLogger()
	// 2.获取配置
	initialize.InitConfig()
	// 3、mysql初始化
	initialize.InitMysql()

	svc := global.ServerConf.ServerInfo
	// 4.随机port测试lb
	debug := initialize.GetEnvInfo("MXSHOP_DEBUG")
	if debug {
		// 使用随机可用端口
		port, err := utils.GetFreePort()
		zap.S().Infof("random port:%d", port)
		if err != nil {
			zap.S().Panicw("service not port", "msg", err.Error())
		}
		svc.Port = port
	}

	// 4、grpc注册
	server := grpc.NewServer()

	proto.RegisterUserServer(server, &handler.UserServer{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", svc.Host, svc.Port))
	if err != nil {
		panic(err)
	}

	// 5、注册健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	// 6、注册consul
	initialize.InitConsul()

	// 7、 服务发现
	go func() {
		zap.S().Infof("gRPC 服务器成功启动在 %s:%d", svc.Host, svc.Port)
		err = server.Serve(listen)
		if err != nil {
			zap.S().Panicw("grpc server start error", zap.Error(err))
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	// 退出时注销服务
	err = global.Consul.Agent().ServiceDeregister(global.ServerConf.ConsulInfo.Id)
	if err != nil {
		zap.S().Errorw("Failed to deregister service", zap.Error(err))
	} else {
		zap.S().Info("Service deregistered successfully")
	}
}
