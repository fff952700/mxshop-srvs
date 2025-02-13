package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"mxshop_srvs/inventory_srv/global"
	"mxshop_srvs/inventory_srv/model"
	"mxshop_srvs/inventory_srv/proto"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (i *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var (
		inv = &model.Inventory{}
	)
	if result := global.DB.Where(&model.Inventory{GoodsId: req.GoodsId}).First(inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "goods inventory not found")
	}
	global.DB.Where("goods_id = ?", req.GoodsId).Updates(&model.Inventory{Stocks: req.Stocks})
	return &emptypb.Empty{}, nil

}

func (i *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inventory = &model.Inventory{}
	if result := global.DB.Where("goods_id = ?", req.GoodsId).First(&inventory); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "inventory not found")
	}
	data := &proto.GoodsInvInfo{
		GoodsId: inventory.GoodsId,
		Stocks:  inventory.Stocks,
	}
	return data, nil

}
func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sell not implemented")
}
func (i *InventoryServer) Rollback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rollback not implemented")
}
