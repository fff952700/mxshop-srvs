package tests

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	"mxshop_srvs/goods_srv/proto"
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

func TestCreateCategory(t *testing.T) {
	InitClient()
	categoryInfo := &proto.CategoryInfoRequest{
		Id:    238013,
		Name:  "testCategory",
		Level: 1,
		IsTab: false,
	}
	rsp, err := GoodsClient.CreateCategory(context.Background(), categoryInfo)
	if err != nil {
		zap.S().Error(err)
	}
	zap.S().Infof("rsp :%v", rsp)
}
