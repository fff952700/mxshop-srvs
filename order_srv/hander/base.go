package hander

import (
	"mxshop_srvs/order_srv/model"
	"mxshop_srvs/order_srv/proto"
)

func (o *OrderServer) Model2InfoResponse(itemInterface interface{}) interface{} {
	switch itemInfo := itemInterface.(type) {
	// 购物车
	case []*model.ShoppingCart:
		var cartItemList []*proto.CartItemResponse
		for _, item := range itemInfo {
			cartItemList = append(cartItemList, &proto.CartItemResponse{
				Id:      item.Id,
				UserId:  item.UserId,
				GoodsId: item.GoodsId,
				Nums:    item.Nums,
				Checked: item.Checked,
			})
		}
		return cartItemList
	default:
		return nil
	}
}
