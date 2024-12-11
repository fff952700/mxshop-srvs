package tests

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/proto"
	"testing"
)

func TestGetCateGoryAll(t *testing.T) {
	InitClient()
	rsp, err := GoodsClient.GetAllCategoryList(context.Background(), &emptypb.Empty{})
	if err != nil {
		zap.S().Error(err)
	}
	zap.S().Info(rsp.JsonData)
}

func TestGetSubCategory(t *testing.T) {
	InitClient()
	rsp, err := GoodsClient.GetSubCategory(context.Background(), &proto.CategoryListRequest{
		Id:    130361,
		Level: 1,
	})
	if err != nil {
		zap.S().Error(err)
	}
	zap.S().Info(rsp.SubCategory)
}
