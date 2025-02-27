package global

import (
	"gorm.io/gorm"

	"github.com/go-redsync/redsync/v4"
	"github.com/hashicorp/consul/api"

	"mxshop_srvs/inventory_srv/config"
)

var (
	ServerConf  = &config.ServerCfg{}
	NacosConf   = &config.NacosConfig{}
	DB          *gorm.DB
	Consul      *api.Client
	RedisClient *redsync.Redsync
)
