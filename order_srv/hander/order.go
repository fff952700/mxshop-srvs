package hander

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

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
	cartItemListResp.Data = data.([]*proto.ShopCartInfoResponse)
	return cartItemListResp, nil
}

func (o *OrderServer) CreateCartItem(ctx context.Context, req *proto.CartItemRequest) (*proto.ShopCartInfoResponse, error) {
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
	return &proto.ShopCartInfoResponse{Id: cartItem.Id}, nil
}

func (o *OrderServer) UpdateCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	// 修改只能修改数量和选中状态 当没有传递值时 proto会使用默认值
	var (
		cartItem  = &model.ShoppingCart{}
		updateMap = make(map[string]interface{})
	)
	if result := global.DB.First(cartItem, req.Id); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	// 当num > 0或者checked 有变化时才修改
	if req.Nums > 0 {
		updateMap["nums"] = req.Nums
	}
	if req.Checked != cartItem.Checked {
		updateMap["checked"] = req.Checked
	}
	if len(updateMap) == 0 {
		return &emptypb.Empty{}, nil
	}
	global.DB.Model(&model.ShoppingCart{}).Where("id =?", req.Id).Updates(updateMap)
	return &emptypb.Empty{}, nil
}

func (o *OrderServer) DeleteCartItem(ctx context.Context, req *proto.CartItemRequest) (*emptypb.Empty, error) {
	var (
		cartItem = &model.ShoppingCart{}
	)

	if result := global.DB.Where(req.Id).First(cartItem); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	// 假删
	//if result := global.DB.Where(req.Id).Updates(&model.ShoppingCart{BaseModel: model.BaseModel{IsDel: true}}); result.RowsAffected != 1 {
	//	return &emptypb.Empty{}, status.Errorf(codes.Internal, "删除失败")
	//}
	global.DB.Delete(cartItem)
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
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{UserId: req.UserId}).Find(&orderList)
	data := o.Model2InfoResponse(orderList)
	return &proto.OrderListResponse{
		Total: int32(total),
		Data:  data.([]*proto.OrderInfoResponse),
	}, nil
}

func (o *OrderServer) OrderDetail(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoDetailResponse, error) {
	// 通用查询。前端查询带有用户id，后台查询只有订单号
	var (
		orderInfo  = &model.OrderInfo{}
		orderGoods []*model.OrderGoods
		orderResp  = &proto.OrderInfoDetailResponse{}
	)
	if result := global.DB.First(orderInfo, req.Id, req.UserId); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	orderResp.OrderInfo = o.Model2InfoResponse(orderInfo).(*proto.OrderInfoResponse)
	if result := global.DB.Where(&model.OrderGoods{OrderSn: orderInfo.OrderSn}).Find(&orderGoods); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单商品不存在")
	}
	orderResp.Goods = o.Model2InfoResponse(orderGoods).([]*proto.OrderItemResponse)
	return orderResp, nil
}

func (o *OrderServer) CreateOrder(ctx context.Context, req *proto.OrderRequest) (*proto.OrderInfoResponse, error) {
	/*
		1. 从购物车获取用户选中的记录
		2. 从商品服务查询商品信息
		3. 从库存服务扣减库存
		4. 订单入库
		5. 购物车删除记录
	*/
	// 查询购物车信息
	var (
		goodsIds      []int32
		total         float32
		sellInfo      = &proto.SellInfo{}
		goodsNumMap   = make(map[int32]int32)
		orderInfoResp = &proto.OrderInfoResponse{}
		shopCartList  []*model.ShoppingCart
		orderGoods    []*model.OrderGoods
		orderSn       = o.GenerateOrderSn()
	)

	if result := global.DB.Where(&model.ShoppingCart{UserId: req.UserId, Checked: true}).Find(&shopCartList); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "购物车记录不存在")
	}

	for _, shop := range shopCartList {
		goodsIds = append(goodsIds, shop.GoodsId)
		goodsNumMap[shop.GoodsId] = shop.Nums
		sellInfo.GoodsInfo = append(sellInfo.GoodsInfo, &proto.GoodsInvInfo{
			GoodsId: shop.GoodsId,
			Stocks:  shop.Nums,
		})
	}
	// 从商品服务查询商品信息
	goodsList, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto.BatchGoodsIdInfo{
		Id: goodsIds,
	})
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "商品信息不存在")
	}
	for _, goods := range goodsList.Data {
		total += goods.ShopPrice * float32(goodsNumMap[goods.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			GoodsId:    goods.Id,
			OrderSn:    orderSn,
			GoodsName:  goods.Name,
			GoodsImage: goods.GoodsFrontImage,
			GoodsPrice: goods.ShopPrice,
			Nums:       goodsNumMap[goods.Id],
		})
	}
	_, err = global.InventorySrvClient.Sell(ctx, sellInfo)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "库存扣减失败")
	}

	// 生成订单
	orderInfoResp.UserId = req.UserId
	// 订单号使用时间戳纳秒级
	orderInfoResp.OrderSn = orderSn
	orderInfoResp.PayType = "not check"
	orderInfoResp.Post = req.Post
	orderInfoResp.Total = total
	orderInfoResp.Address = req.Address
	orderInfoResp.Name = req.Name
	orderInfoResp.Mobile = req.Mobile

	global.DB.Create(&orderInfoResp)

	// 删除购物车记录
	global.DB.Where("user_id =? and checked =?", req.UserId, true).Delete(&model.ShoppingCart{})
	return orderInfoResp, nil
}

// 雪花算法 + 随机数生成订单号
func (o *OrderServer) GenerateOrderSn() string {
	return fmt.Sprintf("%d%d", time.Now().Unix(), rand.Intn(10000))

}
