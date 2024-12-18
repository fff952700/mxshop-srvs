package tests

import (
	"context"
	"go.uber.org/zap"
	"mxshop_srvs/goods_srv/proto"
	"testing"
)

func TestGetBrandList(t *testing.T) {
	rsp, err := GoodsClient.BrandList(context.Background(), &proto.BrandFilterRequest{})
	if err != nil {
		zap.S().Error(err)
	}
	zap.S().Info(rsp.Data)
}
