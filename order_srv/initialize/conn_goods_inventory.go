package initialize

import (
	"context"
	"fmt"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/proto"
)

func init() {
	cfg := global.ServerConf.ConsulInfo
	// 构造对应
	client := []struct {
		name string
		Tag  string
	}{
		{cfg.GoodsName, cfg.GoodsTag},
		{cfg.InventoryName, cfg.InventoryTag},
	}

	for _, c := range client {
		conn, err := grpc.NewClient(fmt.Sprintf("consul://%s:%d/%s?wait=%s&tag=%s", cfg.Host, cfg.Port,
			c.name, "14s", c.Tag),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		)
		fmt.Printf("consul://%s:%d/%s?wait=%s&tag=%s\n", cfg.Host, cfg.Port,
			c.name, "14s", c.Tag)
		if err != nil {
			zap.S().Panicw("Init Client Conn Err", "error", err.Error())
			return
		}
		switch c.name {
		case "goods_srv":
			goodsClient := proto.NewGoodsClient(conn)
			_, err = goodsClient.GetGoodsDetail(context.Background(), &proto.GoodInfoRequest{Id: 421})

			if err != nil {
				zap.S().Panicw("Init UserClient Err", "error", err.Error())
				return
			}
			global.GoodsSrvClient = goodsClient
		case "inventory_srv":
			inventoryClient := proto.NewInventoryClient(conn)
			_, err = inventoryClient.InvDetail(context.Background(), &proto.GoodsInvInfo{GoodsId: 1})

			if err != nil {
				zap.S().Panicw("Init UserClient Err", "error", err.Error())
				return
			}
			global.InventorySrvClient = inventoryClient
		}
	}
}
