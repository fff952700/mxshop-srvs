package tests

import (
	"context"
	"go.uber.org/zap"
	"mxshop_srvs/goods_srv/proto"
	"testing"
)

func TestCategoryBrandList(t *testing.T) {
	InitClient()
	rsp, err := GoodsClient.CategoryBrandList(context.Background(), &proto.CategoryBrandFilterRequest{
		Pages:       1,
		PagePerNums: 10,
	})
	if err != nil {
		zap.S().Error(err)
	}
	zap.S().Info(rsp)
}

func TestGetCategoryBrandList(t *testing.T) {
	InitClient()
	rsp, err := GoodsClient.GetCategoryBrandList(context.Background(), &proto.CategoryInfoRequest{
		Id: 130366,
	})
	if err != nil {
		zap.S().Error(err)
	}
	zap.S().Info(rsp)
}

func TestCreateCategoryBrand(t *testing.T) {
	InitClient()
	rsp, err := GoodsClient.CreateCategoryBrand(context.Background(), &proto.CategoryBrandRequest{
		CategoryId: 130366,
		BrandId:    646,
	})
	if err != nil {
		zap.S().Error(err)
	}
	zap.S().Info(rsp)
}
