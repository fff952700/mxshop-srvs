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
	var (
		inv = &model.Inventory{}
	)
	if result := global.DB.Select("goods_id,stocks").Where(&model.Inventory{GoodsId: req.GoodsId}).First(inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "inventory not found")
	}
	data := &proto.GoodsInvInfo{
		GoodsId: inv.GoodsId,
		Stocks:  inv.Stocks,
	}
	return data, nil

}

// 扣减库存
func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()
	for _, goodsInfo := range req.GoodsInfo {
		for {
			var inv = &model.Inventory{}
			if result := global.DB.Where("goods_id = ?", goodsInfo.GoodsId).First(inv); result.RowsAffected == 0 {
				tx.Rollback()
				return nil, status.Errorf(codes.NotFound, "goods inventory not found")
			}
			if inv.Stocks < goodsInfo.Stocks {
				tx.Rollback()
			}
			inv.Stocks -= goodsInfo.Stocks
			if result := tx.Where("goods_id = ? and version = ?", goodsInfo.GoodsId, inv.Version).Select("stocks", "version").Updates(&model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 1 {
				break
			}

		}
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}

func (i *InventoryServer) Rollback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()
	for _, goodsInfo := range req.GoodsInfo {
		for {
			var inv = &model.Inventory{}
			if result := global.DB.Where("goods_id = ?", goodsInfo.GoodsId).First(inv); result.RowsAffected == 0 {
				tx.Rollback()
				return nil, status.Errorf(codes.NotFound, "goods inventory not found")
			}
			inv.Stocks += goodsInfo.Stocks
			if result := tx.Where("goods_id = ? and version = ?", goodsInfo.GoodsId, inv.Version).Select("stocks", "version").Updates(&model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 1 {
				break
			}
		}
	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
