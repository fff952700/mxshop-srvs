package global

import (
	"github.com/hashicorp/consul/api"
	"gorm.io/gorm"
	"mxshop-srvs/user-srv/config"
)

var (
	ServerConf *config.ServerCfg = &config.ServerCfg{}
	DB         *gorm.DB
	Consul     *api.Client
)
