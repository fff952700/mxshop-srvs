package tests

import (
	"context"
	"go.uber.org/zap"
	"math/rand"
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

func TestBatchGetGoods(t *testing.T) {
	rsp, err := GoodsClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{Id: []int32{421, 422}})
	if err != nil {
		panic(err)
	}
	zap.S().Infof("rsp:%v", rsp.Data)

}

func TestCreateGoods(t *testing.T) {
	rsp, err := GoodsClient.CreateGoods(context.Background(), &proto.CreateGoodsInfo{
		Id:              847,
		Name:            "testCreate",
		GoodsSn:         "",
		Stocks:          0,
		MarketPrice:     rand.Float32(),
		ShopPrice:       rand.Float32(),
		GoodsBrief:      "",
		GoodsDesc:       "",
		ShipFree:        false,
		Images:          nil,
		DescImages:      nil,
		GoodsFrontImage: "",
		IsNew:           false,
		IsHot:           false,
		OnSale:          false,
		CategoryId:      136982,
		BrandId:         614,
	})
	if err != nil {
		panic(err)
	}
	zap.S().Infof("rsp:%v", rsp)

}
