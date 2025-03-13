package hander

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/model"
	"mxshop_srvs/order_srv/proto"
)

type OrderServer struct {
	proto.UnimplementedOrderServer
}

// 购物车相关

func (o *OrderServer) CartItemList(ctx context.Context, req *proto.UserInfo) (*proto.CartItemListResponse, error) {
	var (
		cartItemList     []*model.ShoppingCart
		cartItemListResp = &proto.CartItemListResponse{}
	)
	result := global.DB.Where("user_id = ?", req.UserId).Find(&cartItemList)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}

	cartItemListResp.Total = int32(result.RowsAffected)

	data := o.Model2InfoResponse(cartItemList)
	cartItemListResp.Data = data.([]*proto.CartItemResponse)
	return cartItemListResp, nil
}

func (o *OrderServer) CartItemCreate(ctx context.Context, req *proto.CartItemRequest) (*proto.CartItemResponse, error) {
	var (
		cartItem = &model.ShoppingCart{}
	)
	if result := global.DB.Where("user_id = ? and goods_id = ?", req.UserId, req.GoodsId).Find(cartItem); result.RowsAffected == 0 {
		// 购物车没有该记录
		cartItem.Nums = req.Nums
		cartItem.Checked = req.Checked
		cartItem.GoodsId = req.GoodsId
		cartItem.UserId = req.UserId
	} else {
		cartItem.Nums = req.Nums
		cartItem.Checked = req.Checked
	}
	global.DB.Save(cartItem)
	return &proto.CartItemResponse{Id: cartItem.Id}, nil
}
