package consul

import (
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"

	"mxshop_srvs/order_srv/global"
)

type Register struct {
	Host string
	Port int
}

type RegisterClient interface {
	Register() error
	Deregister() error
}

func NewRegisterClient(host string, port int) RegisterClient {
	return &Register{
		Host: host,
		Port: port,
	}
}

func (r Register) Register() error {
	// 实例化consul对象
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", global.ServerConf.ConsulInfo.Host, global.ServerConf.ConsulInfo.Port)
	client, err := api.NewClient(config)
	if err != nil {
		zap.S().Panicw("[InitConsul] init consul fail")

	}
	global.Consul = client
	// 健康检查
	check := &api.AgentServiceCheck{
		GRPC:     fmt.Sprintf("%s:%d", global.ServerConf.ServerInfo.Host, global.ServerConf.ServerInfo.Port), //
		Timeout:  "5s",                                                                                       // 超时时间
		Interval: "5s",                                                                                       // 运行检查的频率
		// 指定时间后自动注销不健康的服务节点
		DeregisterCriticalServiceAfter: "15s",
	}
	global.ServerConf.ConsulInfo.Id = strconv.FormatInt(time.Now().UnixNano(), 10)
	// 注册consul中的信息 id相同在consul会认为是一个
	registration := &api.AgentServiceRegistration{
		ID:      global.ServerConf.ConsulInfo.Id,   // 服务唯一ID
		Name:    global.ServerConf.ConsulInfo.Name, // 服务名称
		Tags:    global.ServerConf.ConsulInfo.Tag,  // 为服务打标签
		Address: global.ServerConf.ServerInfo.Host,
		Port:    global.ServerConf.ServerInfo.Port,
		Check:   check,
	}
	err = client.Agent().ServiceRegister(registration)
	return err
}

func (r Register) Deregister() error {
	id := global.ServerConf.ConsulInfo.Id
	err := global.Consul.Agent().ServiceDeregister(id)
	return err
}
