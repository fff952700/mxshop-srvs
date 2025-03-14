package hander

import (
	"fmt"
	"gorm.io/gorm"
	"mxshop_srvs/order_srv/model"
	"mxshop_srvs/order_srv/proto"
)

func (o *OrderServer) Model2InfoResponse(itemInterface interface{}) interface{} {
	switch itemInfo := itemInterface.(type) {
	// 购物车
	case []*model.ShoppingCart:
		var cartItemList []*proto.ShopCartInfoResponse
		for _, item := range itemInfo {
			cartItemList = append(cartItemList, &proto.ShopCartInfoResponse{
				Id:      item.Id,
				UserId:  item.UserId,
				GoodsId: item.GoodsId,
				Nums:    item.Nums,
				Checked: item.Checked,
			})
		}
		return cartItemList
	// 订单
	case []*model.OrderInfo:
		var orderList []*proto.OrderInfoResponse
		for _, item := range itemInfo {
			orderList = append(orderList, &proto.OrderInfoResponse{
				Id:      item.Id,
				UserId:  item.UserId,
				OrderSn: item.OrderSn,
				PayType: item.PayType,
				Status:  item.Status,
				Post:    item.Post,
				Address: item.Address,
				Name:    item.SignerName,
				Mobile:  item.SingerMobile,
				PayTime: fmt.Sprintf("%T", item.PayTime),
			})
		}
		return orderList
	case []*model.OrderGoods:
		var orderList []*proto.OrderItemResponse
		for _, item := range itemInfo {
			orderList = append(orderList, &proto.OrderItemResponse{
				GoodsId:    item.GoodsId,
				GoodsName:  item.GoodsName,
				GoodsImage: item.GoodsImage,
				Nums:       item.Nums,
			})
		}
		return orderList
	default:
		return nil
	}
}

func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
