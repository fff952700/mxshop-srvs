package tests

import (
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"mxshop_srvs/goods_srv/proto"
)

var GoodsClient proto.GoodsClient

func init() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	conn, err := grpc.NewClient("192.168.2.150:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Panicw("err conn to server")
	}

	GoodsClient = proto.NewGoodsClient(conn)
}
