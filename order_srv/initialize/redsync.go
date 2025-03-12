package initialize

import (
	"fmt"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"

	"mxshop_srvs/order_srv/global"
)

func init() {
	client := goredislib.NewClient(&goredislib.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConf.RedisInfo.Host, global.ServerConf.RedisInfo.Port),
		DB:   global.ServerConf.RedisInfo.DB,
	})
	rs := redsync.New(goredis.NewPool(client))
	global.RedisClient = rs
}
