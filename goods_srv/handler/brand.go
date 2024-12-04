package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/model"
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
	brandResponses := make([]*proto.BrandInfoResponse, 0)
	for _, brand := range brands {
		brandResponses = append(brandResponses, &proto.BrandInfoResponse{
			Id:   brand.Id,
			Name: brand.Name,
			Logo: brand.Logo,
		})
	}
	var total int64
	global.DB.Model(&model.Brands{}).Count(&total)
	brandListResponse.Data = brandResponses
	brandListResponse.Total = int32(total)
	return &brandListResponse, nil
}

func (g *GoodsServer) CreateBrand(ctx context.Context, req *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	// 查询品牌是否存在
	if result := global.DB.Where("name = ?", req.Name).First(&model.Brands{}); result.RowsAffected == 1 {
		return nil, status.Errorf(codes.InvalidArgument, "商品已存在")
	}
	brand := &model.Brands{
		Name: req.Name,
		Logo: req.Logo,
	}
	global.DB.Create(brand)
	return &proto.BrandInfoResponse{Id: brand.Id}, nil
}

func (g *GoodsServer) DeleteBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.Brands{}, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品不存在")
	}
	return &emptypb.Empty{}, nil

}

func (g *GoodsServer) UpdateBrand(ctx context.Context, req *proto.BrandRequest) (*emptypb.Empty, error) {
	if result := global.DB.Where("name = ?", req.Id).First(&model.Brands{}); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "商品不存在")
	}
	var brand *model.Brands
	if req.Name != "" {
		brand.Name = req.Name
	}
	if req.Logo != "" {
		brand.Name = req.Name
	}
	global.DB.Save(brand)
	return &emptypb.Empty{}, nil
}
