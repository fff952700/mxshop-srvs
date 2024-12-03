package global

import (
	"gorm.io/gorm"

	"github.com/hashicorp/consul/api"

	"mxshop_srvs/user_srv/config"
)

var (
	ServerConf = &config.ServerCfg{}
	NacosConf  = &config.NacosConfig{}
	DB         *gorm.DB
	Consul     *api.Client
)
