package tests

import (
	"context"
	"go.uber.org/zap"
	"mxshop_srvs/goods_srv/proto"
	"testing"
)

func TestGoodsList(t *testing.T) {
	rsp, err := GoodsClient.GoodsList(context.Background(), &proto.GoodsFilterRequest{
		PriceMin: 0,
		PriceMax: 1000,
		//IsHot:       true,
		//IsNew:       true,
		//IsTab:       true,
		TopCategory: 3,
		Pages:       1,
		PagePerNums: 10,
		KeyWords:    "",
		Brand:       0,
	})
	if err != nil {
		panic(err)
	}
	zap.S().Infof("rsp:%v", rsp.Data)

}
