package handler

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"

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

var wg sync.Mutex

// 扣减库存
func (i *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, goods := range req.GoodsInfo {
		// 每次循环创建新实例，避免条件污染
		inv := &model.Inventory{}
		wg.Lock()
		// 明确使用goods_id查询，避免结构体字段污染
		if result := tx.Where("goods_id = ?", goods.GoodsId).First(inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.NotFound, "商品不存在: %d", goods.GoodsId)
		}

		// 检查库存
		if inv.Stocks < goods.Stocks {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足，商品ID: %d", goods.GoodsId)
		}

		// 扣减库存
		inv.Stocks -= goods.Stocks
		if result := tx.Where("goods_id = ?", goods.GoodsId).Updates(&model.Inventory{Stocks: inv.Stocks}); result.Error != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "更新库存失败: %v", result.Error)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "提交事务失败: %v", err)
	}
	wg.Unlock()
	return &emptypb.Empty{}, nil
}

func (i *InventoryServer) Rollback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	var (
		inv = &model.Inventory{}
	)
	// 开启事务保证购物车中的商品要么都成功，要么都失败
	tx := global.DB.Begin()
	for _, goods := range req.GoodsInfo {
		// 判断商品是否存在
		if result := global.DB.Where(&model.Inventory{GoodsId: goods.GoodsId}).First(inv); result.RowsAffected == 0 {
			tx.Rollback()
			return nil, status.Errorf(codes.InvalidArgument, "param err")
		}
		inv.Stocks += goods.Stocks
		// TODO 高并发情况下存在上面条件都通过同时扣减库存的情况。需要引入分布式锁
		tx.Save(inv)

	}
	tx.Commit()
	return &emptypb.Empty{}, nil
}
