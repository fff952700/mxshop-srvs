package hander

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

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
		global.DB.Create(cartItem)
	} else {
		cartItem.Nums += req.Nums
		global.DB.Where("id = ?", cartItem.Id).Updates(&model.ShoppingCart{Nums: cartItem.Nums, Checked: req.Checked})
	}
	return &proto.CartItemResponse{Id: cartItem.Id}, nil
}

func (o *OrderServer) CartItemUpdate(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	// 修改只能修改数量和选中状态 当没有传递值时 proto会使用默认值
	var (
		cartItem = &model.ShoppingCart{}
	)
	if result := global.DB.First(cartItem, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	if req.Nums > 0 {
		cartItem.Nums = req.Nums
	}
	cartItem.Checked = req.Checked
	global.DB.Save(cartItem)
	return &emptypb.Empty{}, nil
}

func (o *OrderServer) CartItemDelete(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	if result := global.DB.Delete(&model.ShoppingCart{}, req.Id); result.RowsAffected != 1 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &emptypb.Empty{}, nil
}

// 订单相关
func (o *OrderServer) OrderList(ctx context.Context, req *proto.OrderFilterRequest) (*proto.OrderListResponse, error) {
	// 通用查询。可能通过前端擦查询或者后台查询所有订单
	var (
		orderList []*model.OrderInfo
		total     int64
	)
	global.DB.Where(&model.OrderInfo{UserId: req.UserId}).Count(&total)
	global.DB.Scopes(Paginate(int(req.Page), int(req.PagePerNums))).Where(&model.OrderInfo{UserId: req.UserId}).Find(&orderList)
	data := o.Model2InfoResponse(orderList)
	return &proto.OrderListResponse{
		Total: int32(total),
		Data:  data.([]*proto.OrderResponse),
	}, nil
}
