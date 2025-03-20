package global

import (
	"gorm.io/gorm"

	"github.com/hashicorp/consul/api"

	"mxshop_srvs/order_srv/config"
	"mxshop_srvs/order_srv/proto"
)

var (
	ServerConf         = &config.ServerCfg{}
	NacosConf          = &config.NacosConfig{}
	DB                 *gorm.DB
	Consul             *api.Client
	InventorySrvClient proto.InventoryClient
	GoodsSrvClient     proto.GoodsClient
)
