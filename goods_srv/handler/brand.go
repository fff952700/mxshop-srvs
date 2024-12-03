package handler

import (
	"context"
	"mxshop_srvs/goods_srv/model"
	"mxshop_srvs/user_srv/global"

	//"google.golang.org/protobuf/types/known/emptypb"
	"mxshop_srvs/goods_srv/proto"
)

func (g *GoodsServer) BrandList(ctx context.Context, req *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	// 实例化返回对象
	brandListResponse := proto.BrandListResponse{}
	// 实例化查询对象
	var brands []model.Brands
	// 分页返回数据 result存在元数据 将数据赋值给brands是通过Find(&brands)完成的
	result := global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Find(&brands)
	if result.Error != nil {
		return &brandListResponse, result.Error
	}
	// 获取总条数
	brandListResponse.Total = int32(result.RowsAffected)
	brandResponses := make([]*proto.BrandInfoResponse, 0)
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.Id,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	return &brandListResponse, nil
}

//func (g *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error){
//
//}
//func (g *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error){
//
//}
//func (g *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error){
//
//}
