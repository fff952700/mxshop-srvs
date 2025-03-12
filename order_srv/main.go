package main

import (
	"fmt"
	"mxshop_srvs/order_srv/hander"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/initialize"
	"mxshop_srvs/order_srv/proto"
	"mxshop_srvs/order_srv/utils"
	"mxshop_srvs/order_srv/utils/register/consul"
)

func main() {
	svc := global.ServerConf.ServerInfo
	// 4.随机port测试lb
	debug := initialize.GetEnvInfo("MXSHOP_DEBUG")
	if !debug {
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

	proto.RegisterOrderServer(server, &hander.OrderServer{})
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", svc.Host, svc.Port))
	if err != nil {
		panic(err)
	}

	// 5、注册健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	client := consul.NewRegisterClient(svc.Host, svc.Port)
	err = client.Register()
	if err != nil {
		zap.S().Panicw("consul register failed", "msg", err.Error())
	}

	// 7、 服务发现
	go func() {
		zap.S().Infof("%s server start sucessfully %s:%d", svc.Name, svc.Host, svc.Port)
		err = server.Serve(listen)
		if err != nil {
			zap.S().Panicf("%s server start error:%T", svc.Name, zap.Error(err))
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	// 退出时注销服务
	err = global.Consul.Agent().ServiceDeregister(global.ServerConf.ConsulInfo.Id)
	if err != nil {
		zap.S().Errorf("%s Failed to deregister service %T", svc.Name, zap.Error(err))
	} else {
		zap.S().Infof("%s Service deregistered successfully", svc.Name)
	}
}
